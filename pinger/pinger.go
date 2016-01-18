package pinger

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/notifier"
)

// Pinger does the HTTP pinging of the sites that are retrieved from the DB.
type Pinger struct {
	Sites      database.Sites
	DB         *sql.DB
	RequestURL URLRequester
	SendEmail  notifier.EmailSender
	SendSms    notifier.SmsSender
	Exit       Exiter
}

// SitesGetter defines a function to get the sites from DB or mock.
type SitesGetter func(db *sql.DB) (database.Sites, error)

// URLRequester defines a function to get thre response and error from http or mock.
type URLRequester func(url string, timeout int) (string, int, time.Duration, error)

// Exiter defines a functio to exit the program to allow exit test scenarios.
type Exiter func(code int)

var stop = make(chan bool)

// InternetAccessError defines errors where the Internet is inaccessible from the server.
const InternetAccessError = "Internet Access Error"

const site1 = "http://www.example.com"
const site2 = "http://www.google.com"

// NewPinger returns a new Pinger object
func NewPinger(db *sql.DB, getSites SitesGetter, requestURL URLRequester,
	exit Exiter, sendEmail notifier.EmailSender, sendSms notifier.SmsSender) *Pinger {
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
		SendSms: sendSms, Exit: exit}
	return &p
}

// Start begins the Pinger service to start pinging
func (p *Pinger) Start() {
	log.Println("Requesting start of pinger...")
	siteCount := 0
	for _, s := range p.Sites {
		//log.Println(s)
		if s.URL != "" {
			go ping(s, p.DB, p.RequestURL, p.SendEmail, p.SendSms)
			siteCount++
		}
	}
	if siteCount == 0 {
		var message = "No active sites set up for pinging in the database!"
		fmt.Println(message)
		log.Println(message)
		p.Exit(1)
	}
}

// Stop stops the Pinger service to end pinging
func (p *Pinger) Stop() {
	log.Println("Requesting stop of pinger...")
	stop <- true
}

// ping does the actual pinging of the site and calls the notifications
func ping(s database.Site, db *sql.DB, requestURL URLRequester,
	sendEmail notifier.EmailSender, sendSms notifier.SmsSender) {
	siteWasUp := true
	var statusChange bool
	var partialDetails string
	var partialSubject string
	for {
		// initialize statusChange to false and only notify on change of siteWasUp status
		statusChange = false
		// Check for a quit signal to stop the pinging
		select {
		case <-stop:
			return
		default:
			if !s.IsActive {
				log.Println(s.Name, "Paused")
				pause(s.PingIntervalSeconds)
				continue
			}
			_, statusCode, responseTime, err := requestURL(s.URL, s.TimeoutSeconds)
			log.Println(s.Name, "Pinged")
			// Record ping informaton.
			p := database.Ping{SiteID: s.SiteID, TimeRequest: time.Now()}
			if err != nil {
				// Check if the error is due to the Internet not being Accessible
				if strings.Contains(err.Error(), InternetAccessError) {
					log.Println(s.Name, "Unable to determine site status -", err)
					pause(s.PingIntervalSeconds)
					continue
				}
				log.Println(s.Name, "Error", err)
				if siteWasUp {
					statusChange = true
					partialSubject = "Site is Down"
					partialDetails = "Site is down, Error is " + err.Error()
				}
				siteWasUp = false
			} else if statusCode != 200 {
				log.Println(s.Name, "Error - HTTP Status Code is", statusCode)
				if siteWasUp {
					statusChange = true
					partialSubject = "Site is Down"
					partialDetails = "Site is down, HTTP Status Code is " + strconv.Itoa(statusCode) + "."
				}
				siteWasUp = false
			} else { // if no errors site is up.
				if !siteWasUp {
					statusChange = true
					partialSubject = "Site is Up"
					partialDetails = fmt.Sprintf("Site is now up, response time was %v.", responseTime)
				}
				siteWasUp = true
			}
			// Save the ping details
			p.Duration = int(responseTime.Nanoseconds() / 1e6)
			p.HTTPStatusCode = statusCode
			p.TimedOut = false
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
			pause(s.PingIntervalSeconds)
		}
	}
}

// pause for the passed number of seconds
func pause(numSeconds int) {
	time.Sleep(time.Duration(numSeconds) * time.Second)
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
			return "", 0, elapsedTime, errors.New(InternetAccessError + ": " + err.Error())
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

// DoExit provides the implementation of the exit function.
func DoExit(flag int) {
	os.Exit(flag)
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
