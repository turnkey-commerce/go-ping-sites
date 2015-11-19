package pinger

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/turnkey-commerce/go-ping-sites/database"
)

// Pinger does the HTTP pinging of the sites that are retrieved from the DB.
type Pinger struct {
	Sites      database.Sites
	DB         *sql.DB
	RequestURL URLRequester
}

// SitesGetter allows to pass a function to get the sites from DB or mock.
type SitesGetter func(db *sql.DB) (database.Sites, error)

// URLRequester allows to pass a function to get thre response and error from http or mock.
type URLRequester func(url string, timeout int) (string, int, error)

var stop = make(chan bool)

// NewPinger returns a new Pinger object
func NewPinger(db *sql.DB, getSites SitesGetter, requestURL URLRequester) *Pinger {
	var sites database.Sites
	var pingerLog *os.File
	var err error
	pingerLog, err = os.Create("pinger.log")
	if err != nil {
		log.Fatal("Error creating pinger log", err)
	}
	log.SetOutput(pingerLog)

	log.Println("Retrieving the initial sites...")
	sites, err = getSites(db)
	if err != nil {
		log.Fatal("Failed to get the sites. ", err)
	}

	for _, s := range sites {
		log.Println("SITE:", s.Name+",", s.URL)
	}

	p := Pinger{Sites: sites, DB: db, RequestURL: requestURL}
	return &p
}

// Start begins the Pinger service to start pinging
func (p *Pinger) Start() {
	log.Println("Requesting start of pinger...")
	siteCount := 0
	for _, s := range p.Sites {
		//log.Println(s)
		if s.URL != "" {
			go ping(s, p.DB, p.RequestURL)
			siteCount++
		}
	}
	if siteCount == 0 {
		log.Println("No active sites set up for pinging.")
	}
}

// Stop stops the Pinger service to end pinging
func (p *Pinger) Stop() {
	log.Println("Requesting stop of pinger...")
	stop <- true
}

// ping does the actual pinging of the site and calls the notifications
func ping(s database.Site, db *sql.DB, requestURL URLRequester) {
	siteWasUp := true
	var notify bool
	var details string
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
					details = "Site is down, Error is " + err.Error()
				}
				siteWasUp = false
			} else if statusCode != 200 {
				log.Println(s.Name, "Error - HTTP Status Code is", statusCode)
				if siteWasUp {
					notify = true
					details = "Site is down, HTTP Status Code is " + strconv.Itoa(statusCode) + "."
				}
				siteWasUp = false
			} else { // if no errors site is up.
				if !siteWasUp {
					notify = true
					details = "Site is now up."
				}
				siteWasUp = true
			}
			if notify {
				log.Println("Will notify status change for", s.Name, details)
			}
			pause(s.PingIntervalSeconds)
		}
	}
}

// pause for the passed number of seconds
func pause(numSeconds int) {
	time.Sleep(time.Duration(numSeconds) * time.Second)
}