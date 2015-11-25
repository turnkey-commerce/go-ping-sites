package pinger

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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
type URLRequester func(url string, timeout int) (string, int, error)

// Exiter defines a functio to exit the program to allow exit test scenarios.
type Exiter func(code int)

var stop = make(chan bool)

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
	var notify bool
	var partialDetails string
	var partialSubject string
	for {
		// initialize notify to false and only notify on change of siteUp status
		notify = false
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
			_, statusCode, err := requestURL(s.URL, s.TimeoutSeconds)
			log.Println(s.Name, "Pinged")
			if err != nil {
				log.Println(s.Name, "Error", err)
				if siteWasUp {
					notify = true
					partialSubject = "Site is Down"
					partialDetails = "Site is down, Error is " + err.Error()
				}
				siteWasUp = false
			} else if statusCode != 200 {
				log.Println(s.Name, "Error - HTTP Status Code is", statusCode)
				if siteWasUp {
					notify = true
					partialSubject = "Site is Down"
					partialDetails = "Site is down, HTTP Status Code is " + strconv.Itoa(statusCode) + "."
				}
				siteWasUp = false
			} else { // if no errors site is up.
				if !siteWasUp {
					notify = true
					partialSubject = "Site is Up"
					partialDetails = "Site is now up."
				}
				siteWasUp = true
			}
			if notify {
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
func RequestURL(url string, timeout int) (string, int, error) {
	to := time.Duration(timeout) * time.Second
	client := http.Client{
		Timeout: to,
	}
	res, err := client.Get(url)
	if err != nil {
		return "", 0, err
	}
	content, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", 0, err
	}
	return string(content), res.StatusCode, nil
}

// GetSites provides the implementation of the SitesGetter type for runtime usage.
func GetSites(db *sql.DB) (database.Sites, error) {
	var sites database.Sites
	err := sites.GetActiveSitesWithContacts(db)
	if err != nil {
		return nil, err
	}
	return sites, nil
}

// DoExit provides the implement of the exit function.
func DoExit(flag int) {
	os.Exit(flag)
}
