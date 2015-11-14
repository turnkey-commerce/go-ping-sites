package pinger

import (
	"database/sql"
	"log"

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
