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

// TestCreateNewSite tests adding a new site in the database and then retrieving it.
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
