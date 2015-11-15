package pinger_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/pinger"
)

// TestNewPinger tests building the pinger object.
func TestNewPinger(t *testing.T) {
	p := pinger.NewPinger(nil, getSites, requestURL)

	if len(p.Sites) != 3 {
		t.Fatal("Incorrect number of sites returned in new pinger.")
	}
}

// TestStartPinger starts up the pinger and then stops after a couple of rounds
func TestStartPinger(t *testing.T) {
	p := pinger.NewPinger(nil, getSites, requestURL)
	p.Start()
	time.Sleep(15 * time.Second)
	p.Stop()
}

func requestURL(url string, timeout int) (string, int, error) {
	if url == "http://www.github.com" {
		return "", 0, errors.New("(Client.Timeout exceeded while awaiting headers)")
	}
	return "Hello", 300, nil
}

func getSites(db *sql.DB) (database.Sites, error) {
	var sites database.Sites
	// Create the first site.
	s1 := database.Site{Name: "Test", IsActive: true, URL: "http://www.google.com",
		PingIntervalSeconds: 5, TimeoutSeconds: 2}
	// Create the second site.
	s2 := database.Site{Name: "Test 2", IsActive: true, URL: "http://www.github.com",
		PingIntervalSeconds: 10, TimeoutSeconds: 5}
	// Create the third site as not active.
	s3 := database.Site{Name: "Test 3", IsActive: false, URL: "http://www.test.com",
		PingIntervalSeconds: 10, TimeoutSeconds: 5}
	// Create first contact
	c1 := database.Contact{Name: "Joe Contact", EmailAddress: "joe@test.com", SmsNumber: "5125551212",
		SmsActive: false, EmailActive: false}
	// Create second contact
	c2 := database.Contact{Name: "Jack Contact", EmailAddress: "jack@test.com", SmsNumber: "5125551213",
		SmsActive: false, EmailActive: false}
	// Add the contacts to the sites
	s1.Contacts = append(s1.Contacts, c1, c2)
	s2.Contacts = append(s2.Contacts, c1)
	s3.Contacts = append(s3.Contacts, c1)

	sites = append(sites, s1, s2, s3)
	return sites, nil
}
