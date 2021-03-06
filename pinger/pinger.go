package pinger

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/notifier"
)

var (
	mu    = &sync.Mutex{}
	site1 = "http://www.example.com"
	site2 = "http://www.google.com"
)

// Pinger does the HTTP pinging of the sites that are retrieved from the DB.
type Pinger struct {
	Sites      database.Sites
	DB         *sql.DB
	RequestURL URLRequester
	SendEmail  notifier.EmailSender
	SendSms    notifier.SmsSender
	getSites   SitesGetter
	wg         sync.WaitGroup
	stopChan   chan struct{}
}

// SitesGetter defines a function to get the sites from DB or mock.
type SitesGetter func(db *sql.DB) (database.Sites, error)

// URLRequester defines a function to get thre response and error from http or mock.
type URLRequester func(url string, timeout int) (string, int, time.Duration, error)

// InternetAccessError defines errors where the Internet is inaccessible from the server.
type InternetAccessError struct {
	msg string
}

func (e InternetAccessError) Error() string {
	return e.msg
}

// NewPinger returns a new Pinger object
func NewPinger(db *sql.DB, getSites SitesGetter, requestURL URLRequester,
	sendEmail notifier.EmailSender, sendSms notifier.SmsSender) *Pinger {
	var sites database.Sites
	var err error

	log.Println("Retrieving the initial sites...")
	sites, err = getSites(db)
	if err != nil {
		// TODO - implement a retry here in case of temporary DB unavailability.
		log.Println("Failed to get the sites. ", err)
	}

	for _, s := range sites {
		log.Println("SITE:", s.Name+",", s.URL)
	}

	p := Pinger{Sites: sites, DB: db, RequestURL: requestURL, SendEmail: sendEmail,
		SendSms: sendSms, getSites: getSites}
	return &p
}

// Start begins the Pinger service to start pinging
func (p *Pinger) Start() {
	log.Println("Requesting start of pingers...")
	siteCount := 0
	p.stopChan = make(chan struct{})
	for _, s := range p.Sites {
		//log.Println(s)
		if s.URL != "" {
			p.wg.Add(1)
			go ping(s, p.DB, p.RequestURL, p.SendEmail, p.SendSms, &p.wg, p.stopChan)
			siteCount++
		}
	}
	if siteCount == 0 {
		var message = "No active sites set up for pinging in the database!"
		fmt.Println(message)
		log.Println(message)
	}
}

// Stop stops the Pinger service by sending stop to all pingers and waits until
// all are stopped via the waitgroup.
func (p *Pinger) Stop() {
	log.Println("Requesting stop of pingers...")
	close(p.stopChan)
	p.wg.Wait()
	// nil out the stopChan to nil so it can be remade when started again.
	p.stopChan = nil
	log.Println("All of the pingers have stopped.")
}

// UpdateSiteSettings stops the pinger, regets the sites for changes in settings,
// and restarts the pinger. There could potentially be race conditions if multiple
// web controllers were trying to update it so a mutex is used to protect it.
func (p *Pinger) UpdateSiteSettings() error {
	// Lock to avoid race conditions since this is usually called from the website.
	mu.Lock()
	defer mu.Unlock()
	log.Println("Updating the site settings due to change...")
	p.Stop()
	// Defer the start in case of error with the get sites.
	defer p.Start()
	sites, err := p.getSites(p.DB)
	if err != nil {
		return err
	}
	p.Sites = sites
	return nil
}

// ping does the actual pinging of the site and calls the notifications
func ping(s database.Site, db *sql.DB, requestURL URLRequester,
	sendEmail notifier.EmailSender, sendSms notifier.SmsSender, wg *sync.WaitGroup, stop chan struct{}) {
	defer wg.Done()
	// Initialize the previous state of site to the database value. On site creation will initialize to true.
	siteWasUp := s.IsSiteUp
	var statusChange bool
	var partialDetails string
	var partialSubject string
	for {
		// initialize statusChange to false and only notify on change of siteWasUp status
		statusChange = false
		// Check for a quit signal to stop the pinging
		select {
		case <-stop:
			log.Println("Stopping ", s.Name)
			return
		case <-time.After(time.Duration(s.PingIntervalSeconds) * time.Second):
			// Do nothing
		}
		if !s.IsActive {
			log.Println(s.Name, "Paused")
			continue
		}
		bodyContent, statusCode, responseTime, err := requestURL(s.URL, s.TimeoutSeconds)
		log.Println(s.Name, "Pinged")
		// Setup ping information for recording.
		p := database.Ping{SiteID: s.SiteID, TimeRequest: time.Now()}
		if err != nil {
			// Check if the error is due to the Internet not being Accessible
			if _, ok := err.(InternetAccessError); ok {
				log.Println(s.Name, "Unable to determine site status -", err)
				continue
			}
			log.Println(s.Name, "Error", err)
			if siteWasUp {
				statusChange = true
				partialSubject = "Site is Down"
				partialDetails = "Site is down, Error is " + err.Error()
			}
			siteWasUp = false

		} else if statusCode < 200 || statusCode > 299 { // Check if the status code is in the 2xx range.
			log.Println(s.Name, "Error - HTTP Status Code is", statusCode)
			if siteWasUp {
				statusChange = true
				partialSubject = "Site is Down"
				partialDetails = "Site is down, HTTP Status Code is " + strconv.Itoa(statusCode) + "."
			}
			siteWasUp = false

		} else {
			siteUp := true
			// if the site settings require check the content.
			if siteUp && s.ContentExpected != "" && !strings.Contains(bodyContent, s.ContentExpected) {
				siteUp = false
				log.Println(s.Name, "Error - required body content missing: ", s.ContentExpected)
				if siteWasUp {
					statusChange = true
					partialSubject = "Site is Down"
					partialDetails = "Site is Down, required body content missing: " + s.ContentExpected + "."
				}
			}
			if siteUp && s.ContentUnexpected != "" && strings.Contains(bodyContent, s.ContentUnexpected) {
				siteUp = false
				log.Println(s.Name, "Error - body content content has excluded content: ", s.ContentUnexpected)
				if siteWasUp {
					statusChange = true
					partialSubject = "Site is Down"
					partialDetails = "Site is Down, body content content has excluded content: " + s.ContentUnexpected + "."
				}
			}
			if siteUp && !siteWasUp {
				statusChange = true
				partialSubject = "Site is Up"
				partialDetails = fmt.Sprintf("Site is now up, response time was %v.", responseTime)
				siteWasUp = true
			}
			siteWasUp = siteUp
		}
		// Save the ping details
		p.Duration = int(responseTime.Nanoseconds() / 1e6)
		p.HTTPStatusCode = statusCode
		p.SiteDown = !siteWasUp
		// Save ping to db.
		err = p.CreatePing(db)
		if err != nil {
			log.Println("Error saving to ping to db:", err)
		}
		// Do the notifications if applicable
		if statusChange {
			// Update the site Status
			err = s.UpdateSiteStatus(db, siteWasUp)
			if err != nil {
				log.Println("Error updating site status:", err)
			}
			// Do the notifications if applicable
			subject := s.Name + ": " + partialSubject
			details := s.Name + " at " + s.URL + ": " + partialDetails
			log.Println("Will notify status change for", s.Name+":", details)

			n := notifier.NewNotifier(s, details, subject, sendEmail, sendSms)
			n.Notify()
		}
	}
}

// RequestURL provides the implementation of the URLRequester type for runtime usage.
func RequestURL(url string, timeout int) (string, int, time.Duration, error) {
	to := time.Duration(timeout) * time.Second
	client := http.Client{
		Timeout: to,
	}
	// Record the timing of the request by diff from the initial time.
	timeStart := time.Now()
	// Do the get request.
	res, err := client.Get(url)
	elapsedTime := round(time.Since(timeStart), time.Millisecond)
	if err != nil {
		// If there's an error need to determine if it could be a local networking error
		// by checking a couple of highly available sites.
		if !isInternetAccessible(site1, site2) {
			return "", 0, elapsedTime, InternetAccessError{msg: err.Error()}
		}
		return "", 0, elapsedTime, err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", 0, elapsedTime, err
	}

	return string(content), res.StatusCode, elapsedTime, nil
}

// GetSites provides the implementation of the SitesGetter type for runtime usage.
func GetSites(db *sql.DB) (database.Sites, error) {
	var sites database.Sites
	// Get active sites with contacts.
	err := sites.GetSites(db, true, true)
	if err != nil {
		return nil, err
	}
	return sites, nil
}

// isInternetAccessible checks two highly available sites to check whether the
// oustide Internet is responding and there are no internal network problems.
func isInternetAccessible(testSite1 string, testSite2 string) bool {
	to := time.Duration(5) * time.Second
	client := http.Client{
		Timeout: to,
	}
	_, err1 := client.Get(testSite1)
	if err1 != nil {
		_, err2 := client.Get(testSite2)
		if err2 != nil {
			return false
		}
	}
	return true
}

// round provides a method to round a time duration.
func round(d, r time.Duration) time.Duration {
	if r <= 0 {
		return d
	}
	neg := d < 0
	if neg {
		d = -d
	}
	if m := d % r; m+m < r {
		d = d - m
	} else {
		d = d + r - m
	}
	if neg {
		return -d
	}
	return d
}
