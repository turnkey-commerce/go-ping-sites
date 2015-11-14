package pinger

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/turnkey-commerce/go-ping-sites/database"
)

// Pinger does the HTTP pinging of the sites that are retrieved from the DB.
type Pinger struct {
	Sites database.Sites
	DB    *sql.DB
}

// NewPinger returns a new Pinger object
func NewPinger(db *sql.DB) *Pinger {
	var sites database.Sites
	var err error
	err = sites.GetActiveSitesWithContacts(db)
	if err != nil {
		log.Fatal("Failed to get the sites.", err)
	}

	p := Pinger{Sites: sites, DB: db}
	return &p
}

// Start begins the Pinter service to start pinging
func (p *Pinger) Start() {
	for _, s := range p.Sites {
		//log.Println(s)
		if s.URL != "" {
			go ping(s)
		}
	}
}

// ping does the actual pinging of the site and calls the notifications
func ping(s database.Site) {
	for {
		fmt.Println(s.Name, time.Now())
		time.Sleep(time.Duration(s.TimeoutSeconds) * time.Second)
	}
}
