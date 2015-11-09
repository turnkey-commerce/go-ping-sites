package database_test

import (
	"reflect"
	"testing"

	"github.com/turnkey-commerce/go-ping-sites/database"
)

// TestCreateDb tests the creation of the database.
func TestCreateDb(t *testing.T) {
	db, err := database.InitializeDB()
	defer db.Close()
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}

	errPing := db.Ping()
	if errPing != nil {
		t.Fatal("Failed to ping database:", errPing)
	}
}

// TestCreateAndGetNewSite tests adding a new site in the database and then retrieving it.
func TestCreateAndGetNewSite(t *testing.T) {
	db, err := database.InitializeDB()
	defer db.Close()
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}

	s := database.Site{Name: "Test", IsActive: true, URL: "http://www.google.com", PingIntervalSeconds: 60, TimeoutSeconds: 30}
	siteID, errCreate := s.CreateSite(db)
	if errCreate != nil {
		t.Fatal("Failed to create new site:", errCreate)
	}

	// Autonumber should be 1 on the first create.
	if siteID != 1 {
		t.Fatal("Expected 1, got ", siteID)
	}

	//Check that the saved site is as Expected
	saved := database.Site{}
	err = saved.GetSite(db, siteID)
	if err != nil {
		t.Fatal("Failed to retrieve new site:", err)
	}
	if !reflect.DeepEqual(s, saved) {
		t.Fatal("New site saved not equal to input:\n", saved, s)
	}
}

// TestCreateUniqueSite tests that the same URL and Site Name can't be entered twice.
func TestCreateUniqueSite(t *testing.T) {
	db, err := database.InitializeDB()
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

// TestCreateAndGetNewContact tests adding a new contact in the database and then retrieving it.
func TestCreateAndGetNewContact(t *testing.T) {
	db, err := database.InitializeDB()
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

	// Create a contact with the side ID
	c := database.Contact{Name: "Joe Contact", EmailAddress: "joe@test.com", SmsNumber: "5125551212",
		SmsActive: false, EmailActive: false}
	contactID, errCreate := c.CreateContact(db, siteID)
	if errCreate != nil {
		t.Fatal("Failed to create new site:", errCreate)
	}

	// Autonumber should be 1 on the first create.
	if contactID != 1 {
		t.Fatal("Expected 1, got ", contactID)
	}
}
