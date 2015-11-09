package database_test

import (
	"database/sql"
	"reflect"
	"testing"

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
	defer db.Close()
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}

	// First create a site to associate the contact.
	s := database.Site{Name: "Test", IsActive: true, URL: "http://www.google.com", PingIntervalSeconds: 60, TimeoutSeconds: 30}
	siteID, errCreate := s.CreateSite(db)
	if errCreate != nil {
		t.Fatal("Failed to create new site:", errCreate)
	}

	// siteID should be 1 on the first create.
	if siteID != 1 {
		t.Fatal("Expected 1, got ", siteID)
	}

	//Get the saved site
	var saved database.Site
	err = saved.GetSite(db, siteID)
	if err != nil {
		t.Fatal("Failed to retrieve new site:", err)
	}
	//Verify the saved site is same as the input.
	if !reflect.DeepEqual(s, saved) {
		t.Fatal("New site saved not equal to input:\n", saved, s)
	}

	// Create first contact with the site ID
	c := database.Contact{Name: "Joe Contact", EmailAddress: "joe@test.com", SmsNumber: "5125551212",
		SmsActive: false, EmailActive: false}
	contactID, err := c.CreateContact(db, siteID)
	if err != nil {
		t.Fatal("Failed to create new site:", errCreate)
	}

	// contactID should be 1 on the first create.
	if contactID != 1 {
		t.Fatal("Expected 1, got ", contactID)
	}

	// Create second contact with the site ID
	c2 := database.Contact{Name: "Jill Contact", EmailAddress: "jill@test.com", SmsNumber: "5125551213",
		SmsActive: false, EmailActive: false}
	_, errCreate2 := c2.CreateContact(db, siteID)
	if errCreate2 != nil {
		t.Fatal("Failed to create new site:", errCreate2)
	}

	//Get the saved site
	err = saved.GetSite(db, siteID)
	if err != nil {
		t.Fatal("Failed to retrieve new site:", err)
	}

	// Verify the first contact was Loaded the same
	if !reflect.DeepEqual(c, saved.Contacts[0]) {
		t.Fatal("New contact saved not equal to input:\n", saved.Contacts[0], c)
	}
	// Verify the second contact was Loaded the same
	if !reflect.DeepEqual(c2, saved.Contacts[1]) {
		t.Fatal("New contact saved not equal to input:\n", saved.Contacts[1], c2)
	}
}

// TestCreateUniqueSite tests that the same URL and Site Name can't be entered twice.
func TestCreateUniqueSite(t *testing.T) {
	db, err := initializeTest()
	defer db.Close()
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}

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
