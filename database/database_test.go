package database_test

import (
	"math"
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
	// Get all of the active sites
	err = sites.GetSites(db, true, false)
	if err != nil {
		t.Error("Failed to get all the sites.", err)
	}

	// Verify that there are  two active sites loaded.
	if len(sites) != 2 {
		t.Error("There should be two active sites loaded.")
	}

	// Verify that GetFirstPing doesn't throw an error with empty Pings
	firstPing, err := sites[0].GetFirstPing(db)
	if err != nil {
		t.Error("GetFirstPing shouldn't throw error if empty pings: ", err)
	}

	zeroTime := time.Time{}
	if firstPing != zeroTime {
		t.Error("GetFirstPing should return a zero time for an empty ping table, but returned: ", err)
	}
}

// TestCreateSiteAndContacts tests creating a site and adding new contacts
// in the database and then retrieving it.
func TestCreateSiteAndContacts(t *testing.T) {
	db, err := database.InitializeTestDB("")
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}
	defer db.Close()

	// First create a site to associate with the contacts.
	// Note: SiteID is ignored for create but is used in the test comparison
	s := database.Site{SiteID: 1, Name: "Test", IsActive: true, URL: "http://www.google.com",
		PingIntervalSeconds: 60, TimeoutSeconds: 30, IsSiteUp: true, ContentExpected: "Expected Content",
		ContentUnexpected: "Unexpected Content"}
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
	if !database.CompareSites(site, s) {
		t.Error("New site saved not equal to input:\n", site, s)
	}

	//Update the saved site
	sUpdate := database.Site{SiteID: 1, Name: "Test Update", IsActive: false,
		URL: "http://www.example.com", PingIntervalSeconds: 30, TimeoutSeconds: 15,
		ContentExpected: "Updated Content", ContentUnexpected: "Updated Unexpected",
		IsSiteUp: true,
	}
	site.Name = sUpdate.Name
	site.URL = sUpdate.URL
	site.IsActive = sUpdate.IsActive
	site.PingIntervalSeconds = sUpdate.PingIntervalSeconds
	site.TimeoutSeconds = sUpdate.TimeoutSeconds
	site.ContentExpected = sUpdate.ContentExpected
	site.ContentUnexpected = sUpdate.ContentUnexpected
	site.IsSiteUp = sUpdate.IsSiteUp
	err = site.UpdateSite(db)
	if err != nil {
		t.Fatal("Failed to update site:", err)
	}

	//Get the updated site
	var siteUpdated database.Site
	err = siteUpdated.GetSite(db, s.SiteID)
	if err != nil {
		t.Fatal("Failed to retrieve updated site:", err)
	}
	//Verify the saved site is same as the input.
	if !database.CompareSites(siteUpdated, sUpdate) {
		t.Error("Updated site saved not equal to input:\n", siteUpdated, sUpdate)
	}

	// Create first contact - ContactID is for referencing the contact get test
	c := database.Contact{Name: "Joe Contact", EmailAddress: "joe@test.com", SmsNumber: "5125551212",
		SmsActive: false, EmailActive: false, ContactID: 1}
	err = c.CreateContact(db)
	if err != nil {
		t.Fatal("Failed to create new contact:", err)
	}
	// Associate to the site ID
	err = site.AddContactToSite(db, c.ContactID)
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
	// Associate the contact to the site
	err = site.AddContactToSite(db, c2.ContactID)
	if err != nil {
		t.Error("Failed to associate contact2 with site:", err)
	}

	//Get the saved site contacts
	err = site.GetSiteContacts(db, site.SiteID)
	if err != nil {
		t.Error("Failed to retrieve site contacts:", err)
	}

	// Verify the first contact was Loaded as the last in list by sort order
	if !reflect.DeepEqual(c, site.Contacts[1]) {
		t.Error("New contact saved not equal to input:\n", site.Contacts[1], c)
	}
	// Verify the second contact was Loaded as the first in list by sort order
	if !reflect.DeepEqual(c2, site.Contacts[0]) {
		t.Error("New contact saved not equal to input:\n", site.Contacts[0], c2)
	}

	// Remove second contact from site.
	err = site.RemoveContactFromSite(db, c2.ContactID)
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

	// Get the first contact via the GetContact method
	c1Get := database.Contact{}
	err = c1Get.GetContact(db, c.ContactID)
	if err != nil {
		t.Error("Failed to retrieve the first contact.")
	}

	// Verify the first contact was retrieved OK
	if !reflect.DeepEqual(c, c1Get) {
		t.Error("Retrieved contact saved not equal to input:\n", c1Get, c)
	}

	// Update the first contact.
	c1Update := database.Contact{Name: "Jane Contact", EmailAddress: "jane@test.com", SmsNumber: "5125551313",
		SmsActive: true, EmailActive: true, ContactID: 1}
	c1Get.Name = c1Update.Name
	c1Get.EmailAddress = c1Update.EmailAddress
	c1Get.SmsNumber = c1Update.SmsNumber
	c1Get.EmailActive = c1Update.EmailActive
	c1Get.SmsActive = c1Update.SmsActive
	err = c1Get.UpdateContact(db)
	if err != nil {
		t.Error("Failed to update the first contact.")
	}

	// Get the first contact again after update
	c1Get2 := database.Contact{}
	err = c1Get2.GetContact(db, c1Update.ContactID)
	if err != nil {
		t.Error("Failed to retrieve the first contact.")
	}

	// Verify the first contact was retrieved OK
	if !reflect.DeepEqual(c1Update, c1Get2) {
		t.Error("Retrieved updated contact saved not equal to input:\n", c1Get2, c1Update)
	}
}

// TestDeleteContacts tests creating a site and two contacts
// in the database and then deleting one of the contacts.
func TestDeleteContact(t *testing.T) {
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

	// Create first contact - ContactID is for referencing the contact get test
	c := database.Contact{Name: "Joe Contact", EmailAddress: "joe@test.com", SmsNumber: "5125551212",
		SmsActive: false, EmailActive: false, ContactID: 1}
	err = c.CreateContact(db)
	if err != nil {
		t.Fatal("Failed to create new contact:", err)
	}
	// Associate to the site ID
	err = s.AddContactToSite(db, c.ContactID)
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
	// Associate the contact to the site
	err = s.AddContactToSite(db, c2.ContactID)
	if err != nil {
		t.Error("Failed to associate contact2 with site:", err)
	}

	err = s.GetSiteContacts(db, s.SiteID)
	if err != nil {
		t.Error("Failed to retrieve site contacts:", err)
	}
	if len(s.Contacts) != 2 {
		t.Error("There should two contacts before deletion.")
	}

	// Delete the second contact
	err = c2.DeleteContact(db)
	if err != nil {
		t.Fatal("Failed to delete contact 2:", err)
	}

	// Verify that it was deleted OK and not associated with the site, and
	// that contact1 is still there.
	err = s.GetSiteContacts(db, s.SiteID)
	if err != nil {
		t.Error("Failed to retrieve site contacts:", err)
	}

	if len(s.Contacts) != 1 {
		t.Error("There should only be one contact for the site after deletion.")
	}
	if !reflect.DeepEqual(c, s.Contacts[0]) {
		t.Error("Remaining contact not equal to input:\n", s.Contacts[0], c)
	}

	// Also verify that the contacts are correct.
	var contacts database.Contacts
	err = contacts.GetContacts(db)
	if err != nil {
		t.Fatal("Failed to get all contacts.", err)
	}

	if len(contacts) != 1 {
		t.Error("There should only be one contact in the DB after deletion.")
	}

	if contacts[0].SiteCount != 1 {
		t.Error("Site count should be 1.")
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

// TestUpdateSiteStatus tests updating the up/down status of the site.
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

	// Update the first ping time of the site.
	firstPingTime := time.Date(2015, time.November, 10, 23, 22, 22, 00, time.UTC)
	err = s.UpdateSiteFirstPing(db, firstPingTime)
	if err != nil {
		t.Fatal("Failed to update first ping time:", err)
	}

	err = updatedSite.GetSite(db, s.SiteID)
	if err != nil {
		t.Fatal("Failed to retrieve updated site:", err)
	}

	if updatedSite.IsSiteUp != true {
		t.Errorf("Site status should be up.")
	}

	if updatedSite.FirstPing != firstPingTime {
		t.Errorf("Site first ping time %s does not match input %s.", updatedSite.FirstPing, firstPingTime)
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
		SmsActive: false, EmailActive: false, SiteCount: 0}
	err = c.CreateContact(db)
	if err != nil {
		t.Fatal("Failed to create new contact:", err)
	}

	// Create second contact
	c2 := database.Contact{Name: "Jack Contact", EmailAddress: "jack@test.com", SmsNumber: "5125551213",
		SmsActive: false, EmailActive: false, SiteCount: 0}
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
		Duration: 280, HTTPStatusCode: 200, SiteDown: false}
	err = p1.CreatePing(db)
	if err != nil {
		t.Fatal("Failed to create new ping:", err)
	}

	// Create a second ping result
	p2 := database.Ping{SiteID: s.SiteID, TimeRequest: time.Date(2015, time.November, 10, 23, 22, 20, 00, time.UTC),
		Duration: 290, HTTPStatusCode: 200, SiteDown: true}
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
		t.Error("First saved ping not equal to input:\n", saved.Pings[1], p1)
	}

	// Verify the second ping was Loaded with proper attributes and sorted first.
	if !reflect.DeepEqual(p2, saved.Pings[0]) {
		t.Error("Second saved ping not equal to input:\n", saved.Pings[0], p2)
	}

	// Verify that the site reflects the last ping time.
	err = saved.GetSite(db, s.SiteID)
	if err != nil {
		t.Fatal("Failed to retrieve site:", err)
	}

	if saved.LastPing != p2.TimeRequest {
		t.Error("Last Ping on site does not match input:\n", saved.LastPing, p1.TimeRequest)
	}

	//Get the first ping for the site.
	firstping, err := s.GetFirstPing(db)
	if err != nil {
		t.Fatal("Failed to retrieve first ping for the site:", err)
	}
	if firstping != p2.TimeRequest {
		t.Error("First Ping on site does not match input:\n", firstping, p2.TimeRequest)
	}

	// Create a third ping with conflicting times should error.
	p3 := database.Ping{SiteID: s.SiteID, TimeRequest: time.Date(2015, time.November, 10, 23, 22, 20, 00, time.UTC),
		Duration: 300, HTTPStatusCode: 200, SiteDown: false}
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
		PingIntervalSeconds: 60, TimeoutSeconds: 30, ContentExpected: "Expected 1",
		ContentUnexpected: "Unexpected 1"}
	err = s1.CreateSite(db)
	if err != nil {
		t.Fatal("Failed to create first site:", err)
	}

	// Create the second site.
	s2 := database.Site{Name: "Test 2", IsActive: true, URL: "http://www.test.com",
		PingIntervalSeconds: 60, TimeoutSeconds: 30,
		LastStatusChange: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		LastPing:         time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		ContentExpected:  "Expected 2", ContentUnexpected: "Unexpected 2"}
	err = s2.CreateSite(db)
	if err != nil {
		t.Fatal("Failed to create second site:", err)
	}

	// Create a third site that is marked inactive.
	s3 := database.Site{Name: "Test 3", IsActive: false, URL: "http://www.test3.com",
		PingIntervalSeconds: 60, TimeoutSeconds: 30, ContentExpected: "Expected 3",
		ContentUnexpected: "Unexpected 3"}
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
	err = s1.AddContactToSite(db, c1.ContactID)
	if err != nil {
		t.Fatal("Failed to associate contact 1 with first site:", err)
	}
	err = s2.AddContactToSite(db, c1.ContactID)
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
	err = s1.AddContactToSite(db, c2.ContactID)
	if err != nil {
		t.Fatal("Failed to associate contact 1 with first site:", err)
	}

	var sites database.Sites
	// Get active sites with contacts
	err = sites.GetSites(db, true, true)
	if err != nil {
		t.Fatal("Failed to get all the sites.", err)
	}

	// Verify that there are only two active sites.
	if len(sites) != 2 {
		t.Fatal("There should only be two active sites loaded.")
	}

	// Verify the first site was Loaded with proper attributes.
	if !database.CompareSites(s1, sites[0]) {
		t.Fatal("First saved site not equal to input:\n", sites[0], s1)
	}

	// Verify the second site was Loaded with proper attributes.
	if !database.CompareSites(s2, sites[1]) {
		t.Fatal("Second saved site not equal to input:\n", sites[1], s2)
	}

	// Verify the first contact was Loaded with proper attributes and sorted last.
	if !reflect.DeepEqual(c1, sites[0].Contacts[1]) {
		t.Error("Second saved contact not equal to input:\n", sites[0].Contacts[1], c1)
	}
	// Verify the second contact was loaded with the proper attributes and sorted first.
	if !reflect.DeepEqual(c2, sites[0].Contacts[0]) {
		t.Error("First saved contact not equal to input:\n", sites[0].Contacts[0], c2)
	}
	// Verify the first contact was loaded to the second site.
	if !reflect.DeepEqual(c1, sites[1].Contacts[0]) {
		t.Error("Second saved contact not equal to input:\n", sites[1].Contacts[0], c1)
	}

	// Verify that the first contact can get both related sites
	err = c1.GetContactSites(db)
	if err != nil {
		t.Error("Error getting the sites for the first contact.")
	}
	if s1.URL != c1.Sites[0].URL {
		t.Error("First contact's first site not as expected:\n", c1.Sites[0].URL, s1.URL)
	}
	if s2.URL != c1.Sites[1].URL {
		t.Error("First contact's second site not as expected:\n", c1.Sites[1].URL, s2.URL)
	}
	if len(c1.Sites) != 2 {
		t.Error("First contact should have two associated site.")
	}

	// Verify that the second contact can get the only related sites
	err = c2.GetContactSites(db)
	if err != nil {
		t.Error("Error getting the site for the second contact.")
	}
	if s1.URL != c2.Sites[0].URL {
		t.Error("Second contact's first site not as expected:\n", c2.Sites[0].URL, s1.URL)
	}
	if len(c2.Sites) != 1 {
		t.Error("Second contact should only have one associated site.")
	}

	// Test for just the active sites without the contacts
	var sitesNoContacts database.Sites
	err = sitesNoContacts.GetSites(db, true, false)
	if err != nil {
		t.Fatal("Failed to get all the sites.", err)
	}

	// Verify the first site was Loaded with proper attributes and no contacts.
	if !database.CompareSites(s1, sitesNoContacts[0]) {
		t.Error("First saved site not equal to GetActiveSites results:\n", sitesNoContacts[0], s1)
	}

	// Verify the second site was Loaded with proper attributes and no contacts.
	if !database.CompareSites(s2, sitesNoContacts[1]) {
		t.Error("Second saved site not equal to GetActiveSites results:\n", sitesNoContacts[1], s2)
	}

	// Test for all of the sites without the contacts
	var allSitesNoContacts database.Sites
	err = allSitesNoContacts.GetSites(db, false, false)
	if err != nil {
		t.Fatal("Failed to get all of the sites.", err)
	}

	// Verify that there are 3 total sites.
	if len(allSitesNoContacts) != 3 {
		t.Error("There should be three total sites loaded.")
	}

}

func TestReportYear(t *testing.T) {
	db, err := database.InitializeReportDB()
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}
	defer db.Close()
	years, err := database.GetReportYears(db)
	if err != nil {
		t.Fatal("Failed to get report years:", err)
	}
	if years[0] != 2017 {
		t.Errorf("Report Year should be 2016, got %d", years[0])
	}
	if years[1] != 2016 {
		t.Errorf("Report Year should be 2017, got %d", years[1])
	}
}

// TestReports verifies the reading of the report from the DB.
func TestReport(t *testing.T) {
	db, err := database.InitializeReportDB()
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}
	defer db.Close()
	report, err := database.GetYTDReports(db, 2016)
	if err != nil {
		t.Fatal("Failed to get report:", err)
	}
	site := "Example.com"
	if report[site][0].PingsUp != 2875 {
		t.Errorf("PingsUp should be 2875, got %d", report[site][0].PingsUp)
	}
	if report[site][0].PingsDown != 0 {
		t.Errorf("PingsDown should be 0, got %d", report[site][0].PingsDown)
	}
	if math.Abs(report[site][0].AvgResponse-37.397565) > .00001 {
		t.Errorf("AvgResponse should be 37.397565, got %f", report[site][0].AvgResponse)
	}
}
