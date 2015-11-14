package pinger_test

import (
	"database/sql"
	"testing"

	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/pinger"
)

// TestNewPinger tests building the pinger object.
func TestNewPinger(t *testing.T) {
	var err error
	db, err := database.InitializeTestDB()
	defer db.Close()
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}
	createTestSites(db, t)

	p := pinger.NewPinger(db)

	if len(p.Sites) != 2 {
		t.Fatal("Incorrect number of sites returned in new pinger.")
	}
}

func createTestSites(db *sql.DB, t *testing.T) {
	var err error
	// Create the first site.
	s1 := database.Site{Name: "Test", IsActive: true, URL: "http://www.google.com",
		PingIntervalSeconds: 60, TimeoutSeconds: 30}
	err = s1.CreateSite(db)
	if err != nil {
		t.Fatal("Failed to create first site:", err)
	}

	// Create the second site.
	s2 := database.Site{Name: "Test 2", IsActive: true, URL: "http://www.test.com",
		PingIntervalSeconds: 60, TimeoutSeconds: 30}
	err = s2.CreateSite(db)
	if err != nil {
		t.Fatal("Failed to create second site:", err)
	}

	// Create a third site that is marked inactive.
	s3 := database.Site{Name: "Test 3", IsActive: false, URL: "http://www.test3.com",
		PingIntervalSeconds: 60, TimeoutSeconds: 30}
	err = s3.CreateSite(db)
	if err != nil {
		t.Fatal("Failed to create third site:", err)
	}

	// Create first contact
	c1 := database.Contact{Name: "Joe Contact", EmailAddress: "joe@test.com", SmsNumber: "5125551212",
		SmsActive: false, EmailActive: false}
	err = c1.CreateContact(db)
	if err != nil {
		t.Fatal("Failed to create new contact:", err)
	}
	// Associate to the first and second site ID
	err = c1.AddContactToSite(db, s1.SiteID)
	if err != nil {
		t.Fatal("Failed to associate contact 1 with first site:", err)
	}
	err = c1.AddContactToSite(db, s2.SiteID)
	if err != nil {
		t.Fatal("Failed to associate contact 1 with second site:", err)
	}

	// Create second contact
	c2 := database.Contact{Name: "Jack Contact", EmailAddress: "jack@test.com", SmsNumber: "5125551213",
		SmsActive: false, EmailActive: false}
	err = c2.CreateContact(db)
	if err != nil {
		t.Fatal("Failed to create new contact:", err)
	}
	// Associate only to the first site
	err = c2.AddContactToSite(db, s1.SiteID)
	if err != nil {
		t.Fatal("Failed to associate contact 1 with first site:", err)
	}
}
