package pinger

import (
	"database/sql"
	"log"
	"time"

	"github.com/turnkey-commerce/go-ping-sites/database"
)

// Pinger does the HTTP pinging of the sites that are retrieved from the DB.
type Pinger struct {
	Sites database.Sites
	DB    *sql.DB
}

// SitesGetter allows to pass a function to get the site from DB or mocks
type SitesGetter func(db *sql.DB) (database.Sites, error)

var stop = make(chan bool)

// NewPinger returns a new Pinger object
func NewPinger(db *sql.DB, getSites SitesGetter) *Pinger {
	var sites database.Sites
	var err error
	sites, err = getSites(db)
	if err != nil {
		log.Fatal("Failed to get the sites. ", err)
	}

	p := Pinger{Sites: sites, DB: db}
	return &p
}

// Start begins the Pinger service to start pinging
func (p *Pinger) Start() {
	log.Println("Requesting start of pinger...")
	for _, s := range p.Sites {
		//log.Println(s)
		if s.URL != "" {
			go ping(s, p.DB)
		}
	}
}

// Stop stops the Pinger service to end pinging
func (p *Pinger) Stop() {
	log.Println("Requesting stop of pinger...")
	stop <- true
}

// ping does the actual pinging of the site and calls the notifications
func ping(s database.Site, db *sql.DB) {
	for {
		// Check for a quit signal to stop the pinging
		select {
		case <-stop:
			return
		default:
			if !s.IsActive {
				pause(s.TimeoutSeconds)
				log.Println(s.Name, "Paused")
				continue
			}
			log.Println(s.Name, "Pinged")
			pause(s.TimeoutSeconds)
		}
	}
}

// pause for the passed number of seconds
func pause(numSeconds int) {
	time.Sleep(time.Duration(numSeconds) * time.Second)
}
