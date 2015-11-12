package database_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/turnkey-commerce/go-ping-sites/database"
)

const testDb string = "./test.db"

func initializeTest() (*sql.DB, error) {
	var db *sql.DB
	err := database.DeleteDb(testDb)
	if err != nil {
		return nil, err
	}
	db, err = database.InitializeDB(testDb)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// TestCreateDb tests the creation of the database.
func TestCreateDb(t *testing.T) {
	db, err := initializeTest()
	if err != nil {
		t.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	errPing := db.Ping()
	if errPing != nil {
		t.Fatal("Failed to ping database:", errPing)
	}
}

// TestCreateSiteAndContacts tests creating a site and adding a new contacts
// in the database and then retrieving it.
func TestCreateSiteAndContacts(t *testing.T) {
	db, err := initializeTest()
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}
	defer db.Close()

	// First create a site to associate with the contacts.
	// Note: SiteID is ignored for create but is used in the test comparison
	s := database.Site{SiteID: 1, Name: "Test", IsActive: true, URL: "http://www.google.com", PingIntervalSeconds: 60, TimeoutSeconds: 30}
	siteID, errCreate := s.CreateSite(db)
	if errCreate != nil {
		t.Fatal("Failed to create new site:", errCreate)
	}

	// siteID should be 1 on the first create.
	if siteID != 1 {
		t.Fatal("Expected 1, got ", siteID)
	}

	//Get the saved site
	var site database.Site
	err = site.GetSite(db, siteID)
	if err != nil {
		t.Fatal("Failed to retrieve new site:", err)
	}
	//Verify the saved site is same as the input.
	if !reflect.DeepEqual(s, site) {
		t.Fatal("New site saved not equal to input:\n", site, s)
	}

	// Create first contact
	c := database.Contact{Name: "Joe Contact", EmailAddress: "joe@test.com", SmsNumber: "5125551212",
		SmsActive: false, EmailActive: false}
	err = c.CreateContact(db)
	if err != nil {
		t.Fatal("Failed to create new contact:", errCreate)
	}
	// Associate to the site ID
	err = c.AddContactToSite(db, siteID)
	if err != nil {
		t.Fatal("Failed to associate contact with site:", errCreate)
	}

	// Create second contact with the site ID
	c2 := database.Contact{Name: "Jill Contact", EmailAddress: "jill@test.com", SmsNumber: "5125551213",
		SmsActive: false, EmailActive: false}
	err = c2.CreateContact(db)
	if err != nil {
		t.Fatal("Failed to create new site:", err)
	}
	// Associate to the site ID
	err = c2.AddContactToSite(db, siteID)
	if err != nil {
		t.Fatal("Failed to associate contact2 with site:", errCreate)
	}

	//Get the saved site
	err = site.GetSiteContacts(db, siteID)
	if err != nil {
		t.Fatal("Failed to retrieve site contacts:", err)
	}

	// Verify the first contact was Loaded as the last in list by sort order
	if !reflect.DeepEqual(c, site.Contacts[1]) {
		t.Fatal("New contact saved not equal to input:\n", site.Contacts[1], c)
	}
	// Verify the second contact was Loaded as the first in list by sort order
	if !reflect.DeepEqual(c2, site.Contacts[0]) {
		t.Fatal("New contact saved not equal to input:\n", site.Contacts[0], c2)
	}

	// Remove second contact from site.
	err = c2.RemoveContactFromSite(db, siteID)
	if err != nil {
		t.Fatal("Failed to remove contact2 from site:", err)
	}

	//Get the saved site contacts again
	err = site.GetSiteContacts(db, siteID)
	if err != nil {
		t.Fatal("Failed to retrieve site contacts:", err)
	}

	if len(site.Contacts) != 1 {
		t.Fatal("Site should have only one contact after removal")
	}

}

// TestCreateUniqueSite tests that the same URL and Site Name can't be entered twice.
func TestCreateUniqueSite(t *testing.T) {
	db, err := initializeTest()
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}
	defer db.Close()

	s := database.Site{Name: "Test", IsActive: true, URL: "http://www.test.com",
		PingIntervalSeconds: 60, TimeoutSeconds: 30}
	_, errCreate := s.CreateSite(db)
	if errCreate != nil {
		t.Fatal("Failed to create new site:", errCreate)
	}

	//Test where the URL is the same and the Name is different should fail uniqueness constraint.
	s2 := database.Site{Name: "Test2", IsActive: true, URL: "http://www.test.com",
		PingIntervalSeconds: 60, TimeoutSeconds: 30}
	_, errCreate2 := s2.CreateSite(db)
	if errCreate2 == nil {
		t.Fatal("Should throw uniqueness constraint error for URL.")
	}

	//Test where the Name is the same and the URL is different should fail with uniqueness constraint.
	s3 := database.Site{Name: "Test", IsActive: true, URL: "http://www.test.edu",
		PingIntervalSeconds: 60, TimeoutSeconds: 30}
	_, errCreate3 := s3.CreateSite(db)
	if errCreate3 == nil {
		t.Fatal("Should throw uniqueness constraint error for Name.")
	}
}

func TestCreateAndGetUnattachedContacts(t *testing.T) {
	db, err := initializeTest()
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}
	defer db.Close()

	// Create first contact
	c := database.Contact{Name: "Joe Contact", EmailAddress: "joe@test.com", SmsNumber: "5125551212",
		SmsActive: false, EmailActive: false}
	err = c.CreateContact(db)
	if err != nil {
		t.Fatal("Failed to create new contact:", err)
	}

	// Create second contact
	c2 := database.Contact{Name: "Jack Contact", EmailAddress: "jack@test.com", SmsNumber: "5125551213",
		SmsActive: false, EmailActive: false}
	err = c2.CreateContact(db)
	if err != nil {
		t.Fatal("Failed to create new contact:", err)
	}

	// Create third contact with name conflict.
	c3 := database.Contact{Name: "Jack Contact", EmailAddress: "jack@test.com", SmsNumber: "5125551213",
		SmsActive: false, EmailActive: false}
	err = c3.CreateContact(db)
	if err == nil {
		t.Fatal("Conflicting contact should throw error.")
	}

	var contacts database.Contacts
	err = contacts.GetContacts(db)
	if err != nil {
		t.Fatal("Failed to get all contacts.", err)
	}

	fmt.Println(contacts)

	// Verify the first contact was Loaded as the last in list by sort order
	if !reflect.DeepEqual(c, contacts[1]) {
		t.Fatal("New contact saved not equal to input:\n", contacts[1], c)
	}
	// Verify the second contact was Loaded as the first in list by sort order
	if !reflect.DeepEqual(c2, contacts[0]) {
		t.Fatal("New contact saved not equal to input:\n", contacts[0], c2)
	}
}

// TestCreateSiteAndContacts tests creating a site and adding a new contacts
// in the database and then retrieving it.
func TestCreatePings(t *testing.T) {
	db, err := initializeTest()
	defer db.Close()
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}

	// First create a site to associate with the pings.
	s := database.Site{SiteID: 1, Name: "Test", IsActive: true, URL: "http://www.google.com", PingIntervalSeconds: 60, TimeoutSeconds: 30}
	siteID, errCreate := s.CreateSite(db)
	if errCreate != nil {
		t.Fatal("Failed to create new site:", errCreate)
	}

	// Create a ping result
	p1 := database.Ping{SiteID: siteID, TimeRequest: time.Date(2015, time.November, 10, 23, 22, 22, 00, time.UTC), TimeResponse: time.Date(2015, time.November, 10, 23, 22, 25, 00, time.UTC), HTTPStatusCode: 200, TimedOut: false}
	errCreate = p1.CreatePing(db)
	if errCreate != nil {
		t.Fatal("Failed to create new ping:", errCreate)
	}

	// Create a second ping result
	p2 := database.Ping{SiteID: siteID, TimeRequest: time.Date(2015, time.November, 10, 23, 22, 20, 00, time.UTC), TimeResponse: time.Date(2015, time.November, 10, 23, 22, 25, 00, time.UTC), HTTPStatusCode: 200, TimedOut: false}
	errCreate = p2.CreatePing(db)
	if errCreate != nil {
		t.Fatal("Failed to create new ping:", errCreate)
	}

	//Get the saved Ping
	var saved database.Site
	err = saved.GetSitePings(db, siteID, time.Date(2015, time.November, 10, 23, 00, 00, 00, time.UTC),
		time.Date(2015, time.November, 10, 23, 59, 00, 00, time.UTC))
	if err != nil {
		t.Fatal("Failed to retrieve saved pings:", err)
	}

	// Verify the first ping was Loaded with proper attibutes and sorted last.
	if !reflect.DeepEqual(p1, saved.Pings[1]) {
		t.Fatal("First saved ping not equal to input:\n", saved.Pings[1], p1)
	}

	// Verify the second ping was Loaded with proper attributes and sorted first.
	if !reflect.DeepEqual(p2, saved.Pings[0]) {
		t.Fatal("Second saved ping not equal to input:\n", saved.Pings[0], p2)
	}

	// Create a third ping with conflicting times should error.
	p3 := database.Ping{SiteID: siteID, TimeRequest: time.Date(2015, time.November, 10, 23, 22, 20, 00, time.UTC), TimeResponse: time.Date(2015, time.November, 10, 23, 22, 25, 00, time.UTC), HTTPStatusCode: 200, TimedOut: false}
	errCreate = p3.CreatePing(db)
	if errCreate == nil {
		t.Fatal("Conflicting pings should throw error.")
	}
}
