package pinger_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/pinger"
)

// TestNewPinger tests building the pinger object.
func TestNewPinger(t *testing.T) {
	p := pinger.NewPinger(nil, getSites)

	if len(p.Sites) != 2 {
		t.Fatal("Incorrect number of sites returned in new pinger.")
	}
}

// TestStartPinger starts up the pinger and then stops after a couple of rounds
func TestStartPinger(t *testing.T) {
	p := pinger.NewPinger(nil, getSites)
	t.Log("Starting Pinger...")
	p.Start()
	time.Sleep(30 * time.Second)
	t.Log("Stopping Pinger...")
	p.Stop()
}

func getSites(db *sql.DB) (database.Sites, error) {
	var sites database.Sites
	// Create the first site.
	s1 := database.Site{Name: "Test", IsActive: true, URL: "http://www.google.com",
		PingIntervalSeconds: 10, TimeoutSeconds: 5}
	// Create the second site.
	s2 := database.Site{Name: "Test 2", IsActive: true, URL: "http://www.github.com",
		PingIntervalSeconds: 15, TimeoutSeconds: 10}
	// Create first contact
	c1 := database.Contact{Name: "Joe Contact", EmailAddress: "joe@test.com", SmsNumber: "5125551212",
		SmsActive: false, EmailActive: false}
	// Create second contact
	c2 := database.Contact{Name: "Jack Contact", EmailAddress: "jack@test.com", SmsNumber: "5125551213",
		SmsActive: false, EmailActive: false}
	// Add the contacts to the sites
	s1.Contacts = append(s1.Contacts, c1, c2)
	s2.Contacts = append(s2.Contacts, c1)

	sites = append(sites, s1, s2)
	return sites, nil
}
