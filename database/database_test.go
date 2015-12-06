package database_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/turnkey-commerce/go-ping-sites/database"
)

// TestCreateDb tests the creation and initial seeding of the database.
func TestCreateDb(t *testing.T) {
	db, err := database.InitializeTestDB("db-seed.toml")
	if err != nil {
		t.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	errPing := db.Ping()
	if errPing != nil {
		t.Fatal("Failed to ping database:", errPing)
	}

	var sites database.Sites
	err = sites.GetActiveSitesWithContacts(db)
	if err != nil {
		t.Fatal("Failed to get all the sites.", err)
	}

	// Verify that there are  two active sites loaded.
	if len(sites) != 2 {
		t.Fatal("There should be two active sites loaded.")
	}
}

// TestCreateSiteAndContacts tests creating a site and adding a new contacts
// in the database and then retrieving it.
func TestCreateSiteAndContacts(t *testing.T) {
	db, err := database.InitializeTestDB("")
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}
	defer db.Close()

	// First create a site to associate with the contacts.
	// Note: SiteID is ignored for create but is used in the test comparison
	s := database.Site{SiteID: 1, Name: "Test", IsActive: true, URL: "http://www.google.com", PingIntervalSeconds: 60, TimeoutSeconds: 30, IsSiteUp: true}
	err = s.CreateSite(db)
	if err != nil {
		t.Fatal("Failed to create new site:", err)
	}

	// siteID should be 1 on the first create.
	if s.SiteID != 1 {
		t.Fatal("Expected 1, got ", s.SiteID)
	}

	//Get the saved site
	var site database.Site
	err = site.GetSite(db, s.SiteID)
	if err != nil {
		t.Fatal("Failed to retrieve new site:", err)
	}
	//Verify the saved site is same as the input.
	if site.URL != s.URL || site.IsActive != s.IsActive || site.Name != s.Name ||
		site.TimeoutSeconds != s.TimeoutSeconds || site.PingIntervalSeconds != s.PingIntervalSeconds ||
		site.SiteID != s.SiteID {
		t.Error("New site saved not equal to input:\n", site, s)
	}

	// Create first contact
	c := database.Contact{Name: "Joe Contact", EmailAddress: "joe@test.com", SmsNumber: "5125551212",
		SmsActive: false, EmailActive: false}
	err = c.CreateContact(db)
	if err != nil {
		t.Fatal("Failed to create new contact:", err)
	}
	// Associate to the site ID
	err = c.AddContactToSite(db, site.SiteID)
	if err != nil {
		t.Fatal("Failed to associate contact with site:", err)
	}

	// Create second contact
	c2 := database.Contact{Name: "Jill Contact", EmailAddress: "jill@test.com", SmsNumber: "5125551213",
		SmsActive: false, EmailActive: false}
	err = c2.CreateContact(db)
	if err != nil {
		t.Fatal("Failed to create new site:", err)
	}
	// Associate to the site ID
	err = c2.AddContactToSite(db, site.SiteID)
	if err != nil {
		t.Fatal("Failed to associate contact2 with site:", err)
	}

	//Get the saved site
	err = site.GetSiteContacts(db, site.SiteID)
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
	err = c2.RemoveContactFromSite(db, site.SiteID)
	if err != nil {
		t.Fatal("Failed to remove contact2 from site:", err)
	}

	//Get the saved site contacts again
	err = site.GetSiteContacts(db, site.SiteID)
	if err != nil {
		t.Fatal("Failed to retrieve site contacts:", err)
	}

	if len(site.Contacts) != 1 {
		t.Fatal("Site should have only one contact after removal")
	}

}

// TestCreateUniqueSite tests that the same URL and Site Name can't be entered twice.
func TestCreateUniqueSite(t *testing.T) {
	var err error
	db, err := database.InitializeTestDB("")
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}
	defer db.Close()

	s := database.Site{Name: "Test", IsActive: true, URL: "http://www.test.com",
		PingIntervalSeconds: 60, TimeoutSeconds: 30}
	err = s.CreateSite(db)
	if err != nil {
		t.Fatal("Failed to create new site:", err)
	}

	//Test where the URL is the same and the Name is different should fail uniqueness constraint.
	s2 := database.Site{Name: "Test2", IsActive: true, URL: "http://www.test.com",
		PingIntervalSeconds: 60, TimeoutSeconds: 30}
	err = s2.CreateSite(db)
	if err == nil {
		t.Fatal("Should throw uniqueness constraint error for URL.")
	}

	//Test where the Name is the same and the URL is different should fail with uniqueness constraint.
	s3 := database.Site{Name: "Test", IsActive: true, URL: "http://www.test.edu",
		PingIntervalSeconds: 60, TimeoutSeconds: 30}
	err = s3.CreateSite(db)
	if err == nil {
		t.Fatal("Should throw uniqueness constraint error for Name.")
	}
}

// TestCreateSiteAndContacts tests creating a site and adding a new contacts
// in the database and then retrieving it.
func TestUpdateSiteStatus(t *testing.T) {
	db, err := database.InitializeTestDB("")
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}
	defer db.Close()

	// First create a site to update status.
	s := database.Site{SiteID: 1, Name: "Test", IsActive: true, URL: "http://www.google.com", PingIntervalSeconds: 60, TimeoutSeconds: 30}
	err = s.CreateSite(db)
	if err != nil {
		t.Fatal("Failed to create new site:", err)
	}

	// Update the status of the site to down
	err = s.UpdateSiteStatus(db, false)
	if err != nil {
		t.Fatal("Failed to update site status:", err)
	}

	//Get the saved site
	var updatedSite database.Site
	err = updatedSite.GetSite(db, s.SiteID)
	if err != nil {
		t.Fatal("Failed to retrieve updated site:", err)
	}

	if updatedSite.IsSiteUp != false {
		t.Errorf("Site status should be down.")
	}

	// Update the status of the site to up
	err = s.UpdateSiteStatus(db, true)
	if err != nil {
		t.Fatal("Failed to update site status:", err)
	}

	err = updatedSite.GetSite(db, s.SiteID)
	if err != nil {
		t.Fatal("Failed to retrieve updated site:", err)
	}

	if updatedSite.IsSiteUp != true {
		t.Errorf("Site status should be up.")
	}
}

// TestCreateAndGetUnattachedContacts tests the creation of contacts not associated with a site.
func TestCreateAndGetUnattachedContacts(t *testing.T) {
	db, err := database.InitializeTestDB("")
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

	// Verify the first contact was Loaded as the last in list by sort order
	if !reflect.DeepEqual(c, contacts[1]) {
		t.Fatal("New contact saved not equal to input:\n", contacts[1], c)
	}
	// Verify the second contact was Loaded as the first in list by sort order
	if !reflect.DeepEqual(c2, contacts[0]) {
		t.Fatal("New contact saved not equal to input:\n", contacts[0], c2)
	}
}

// TestCreatePings tests creating the ping records for a given site.
func TestCreatePings(t *testing.T) {
	var err error
	db, err := database.InitializeTestDB("")
	defer db.Close()
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}

	// First create a site to associate with the pings.
	s := database.Site{Name: "Test", IsActive: true, URL: "http://www.google.com", PingIntervalSeconds: 60, TimeoutSeconds: 30}
	err = s.CreateSite(db)
	if err != nil {
		t.Fatal("Failed to create new site:", err)
	}

	// Create a ping result
	p1 := database.Ping{SiteID: s.SiteID, TimeRequest: time.Date(2015, time.November, 10, 23, 22, 22, 00, time.UTC),
		Duration: 280, HTTPStatusCode: 200, TimedOut: false}
	err = p1.CreatePing(db)
	if err != nil {
		t.Fatal("Failed to create new ping:", err)
	}

	// Create a second ping result
	p2 := database.Ping{SiteID: s.SiteID, TimeRequest: time.Date(2015, time.November, 10, 23, 22, 20, 00, time.UTC),
		Duration: 290, HTTPStatusCode: 200, TimedOut: false}
	err = p2.CreatePing(db)
	if err != nil {
		t.Fatal("Failed to create new ping:", err)
	}

	//Get the saved Ping
	var saved database.Site
	err = saved.GetSitePings(db, s.SiteID, time.Date(2015, time.November, 10, 23, 00, 00, 00, time.UTC),
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
	p3 := database.Ping{SiteID: s.SiteID, TimeRequest: time.Date(2015, time.November, 10, 23, 22, 20, 00, time.UTC),
		Duration: 300, HTTPStatusCode: 200, TimedOut: false}
	err = p3.CreatePing(db)
	if err == nil {
		t.Fatal("Conflicting pings should throw error.")
	}
}

// TestCreateAndGetMultipleSites tests creating more than one active sites
// with contacts in the database and then retrieving them.
func TestCreateAndGetMultipleSites(t *testing.T) {
	var err error
	db, err := database.InitializeTestDB("")
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}
	defer db.Close()

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

	var sites database.Sites
	err = sites.GetActiveSitesWithContacts(db)
	if err != nil {
		t.Fatal("Failed to get all the sites.", err)
	}

	// Verify that there are only two active sites.
	if len(sites) != 2 {
		t.Fatal("There should only be two active sites loaded.")
	}

	// Verify the first site was Loaded with proper attributes.
	if s1.URL != sites[0].URL || s1.IsActive != sites[0].IsActive ||
		s1.Name != sites[0].Name || s1.PingIntervalSeconds != sites[0].PingIntervalSeconds ||
		s1.TimeoutSeconds != sites[0].TimeoutSeconds || s1.SiteID != sites[0].SiteID {
		t.Fatal("First saved site not equal to input:\n", sites[0], s1)
	}

	// Verify the second site was Loaded with proper attributes.
	if s2.URL != sites[1].URL || s1.IsActive != sites[1].IsActive ||
		s2.Name != sites[1].Name || s2.PingIntervalSeconds != sites[1].PingIntervalSeconds ||
		s2.TimeoutSeconds != sites[1].TimeoutSeconds || s2.SiteID != sites[1].SiteID {
		t.Fatal("Second saved site not equal to input:\n", sites[1], s2)
	}

	// Verify the first contact was Loaded with proper attributes and sorted last.
	if !reflect.DeepEqual(c1, sites[0].Contacts[1]) {
		t.Fatal("Second saved contact not equal to input:\n", sites[0].Contacts[1], c1)
	}
	// Verify the second contact was loaded with the proper attributes and sorted first.
	if !reflect.DeepEqual(c2, sites[0].Contacts[0]) {
		t.Fatal("First saved contact not equal to input:\n", sites[0].Contacts[0], c2)
	}
	// Verify the first contact was loaded to the second site.
	if !reflect.DeepEqual(c1, sites[1].Contacts[0]) {
		t.Fatal("Second saved contact not equal to input:\n", sites[1].Contacts[0], c1)
	}
}
