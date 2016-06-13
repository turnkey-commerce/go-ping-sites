package pinger_test

import (
	"database/sql"
	"reflect"
	"strings"
	"testing"
	"time"

	_ "github.com/erikstmartin/go-testdb"
	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/notifier"
	"github.com/turnkey-commerce/go-ping-sites/pinger"
)

// TestNewPinger tests building the pinger object.
func TestNewPinger(t *testing.T) {
	db, _ := sql.Open("testdb", "")
	pinger.CreatePingerLog("", true)
	p := pinger.NewPinger(db, pinger.GetSitesMock, pinger.RequestURLMock,
		notifier.SendEmailMock, notifier.SendSmsMock)

	if len(p.Sites) != 3 {
		t.Fatal("Incorrect number of sites returned in new pinger.")
	}

	results, err := pinger.GetLogContent()
	if err != nil {
		t.Fatal("Failed to get log results.", err)
	}

	if !strings.Contains(results, "SITE: Test, http://www.google.com") {
		t.Fatal("Failed to load first site.")
	}
	if !strings.Contains(results, "SITE: Test 2, http://www.github.com") {
		t.Fatal("Failed to load second site.")
	}
	if !strings.Contains(results, "SITE: Test 3, http://www.test.com") {
		t.Fatal("Failed to load third site.")
	}
}

// TestStartEmptySitesPinger verifies that proper reporting is done for the case of no active sites.
func TestStartEmptySitesPinger(t *testing.T) {
	db, _ := sql.Open("testdb", "")
	pinger.CreatePingerLog("", true)
	p := pinger.NewPinger(db, pinger.GetEmptySitesMock, pinger.RequestURLMock,
		notifier.SendEmailMock, notifier.SendSmsMock)
	p.Start()

	results, err := pinger.GetLogContent()
	if err != nil {
		t.Fatal("Failed to get log results.", err)
	}

	if !strings.Contains(results, "No active sites set up for pinging in the database!") {
		t.Fatal("Failed to report empty sites.")
	}
}

// TestStartPingerErrorWithGetSites verifies that an error is handled when the get sites returns
// an error.
func TestStartPingerErrorWithGetSites(t *testing.T) {
	db, _ := sql.Open("testdb", "")
	pinger.CreatePingerLog("", true)
	p := pinger.NewPinger(db, pinger.GetSitesErrorMock, pinger.RequestURLMock,
		notifier.SendEmailMock, notifier.SendSmsMock)
	p.Start()

	results, err := pinger.GetLogContent()
	if err != nil {
		t.Fatal("Failed to get log results.", err)
	}

	if !strings.Contains(results, "Timeout accessing the SQL database.") {
		t.Fatal("Failed to report error getting the sites from the DB.")
	}
}

// TestStartAndRestartPinger starts up the pinger and then stops it after 3 seconds
func TestStartAndRestartPinger(t *testing.T) {
	// Fake db for testing.
	db, _ := sql.Open("testdb", "")
	pinger.CreatePingerLog("", true)
	p := pinger.NewPinger(db, pinger.GetSitesMock, pinger.RequestURLMock,
		notifier.SendEmailMock, notifier.SendSmsMock)
	p.Start()
	// Sleep to allow running the tests before stopping.
	time.Sleep(5 * time.Second)
	p.Stop() // Test Restart after stop
	p.Start()
	p.Stop()

	results, err := pinger.GetLogContent()
	if err != nil {
		t.Fatal("Failed to get log results.", err)
	}

	if !strings.Contains(results, "Client.Timeout") {
		t.Errorf("Failed to report timeout error.")
	}
	if !strings.Contains(results, "Test 3 Paused") {
		t.Errorf("Failed to report paused site.")
	}
	if !strings.Contains(results, "Error - HTTP Status Code") {
		t.Fatal("Failed to report bad HTTP Status Code.")
	}
	if !strings.Contains(results, "Sending Notification of Site Contacts about Test 2: Site is Down...") {
		t.Fatal("Failed to report site being down.")
	}
	if !strings.Contains(results, "Will notify status change for Test 2: Test 2 at http://www.github.com: Site is now up, response time was 300ms.") {
		t.Fatal("Failed to report change in notification.")
	}
}

// TestUpdateSiteSettings starts up the pinger and then updates the site settings.
func TestUpdateSiteSettings(t *testing.T) {
	// Fake db for testing.
	db, _ := sql.Open("testdb", "")
	pinger.CreatePingerLog("", true)
	p := pinger.NewPinger(db, pinger.GetSitesMock, pinger.RequestURLMock,
		notifier.SendEmailMock, notifier.SendSmsMock)
	p.Start()
	// Test UpdateSiteSettings
	p.UpdateSiteSettings()
	p.Stop()

	results, err := pinger.GetLogContent()
	if err != nil {
		t.Fatal("Failed to get log results.", err)
	}
	if !strings.Contains(results, "Updating the site settings due to change...") {
		t.Fatal("Failed to launch update site settings.")
	}
}

// TestUpdateSiteSettingsError tests when UpdateSiteSettings returns a DB error.
func TestUpdateSiteSettingsError(t *testing.T) {
	// Fake db for testing.
	db, _ := sql.Open("testdb", "")
	pinger.CreatePingerLog("", true)
	p := pinger.NewPinger(db, pinger.GetSitesErrorMock, pinger.RequestURLMock,
		notifier.SendEmailMock, notifier.SendSmsMock)
	p.Start()
	// Test UpdateSiteSettings with error due to database.
	err := p.UpdateSiteSettings()
	if err == nil {
		t.Error("Failed to report error with getting site settings.")
	}
	p.Stop()

	results, err := pinger.GetLogContent()
	if err != nil {
		t.Fatal("Failed to get log results.", err)
	}

	if !strings.Contains(results, "Timeout accessing the SQL database.") {
		t.Fatal("Failed to report problem updating the sites.")
	}
}

// TestInternetAccessError starts up the pinger and then stops it after 3 seconds
func TestInternetAccessError(t *testing.T) {
	// Fake db for testing.
	db, _ := sql.Open("testdb", "")
	pinger.CreatePingerLog("", true)
	p := pinger.NewPinger(db, pinger.GetSitesMock, pinger.RequestURLBadInternetAccessMock,
		notifier.SendEmailMock, notifier.SendSmsMock)
	p.Start()
	time.Sleep(5 * time.Second)
	p.Stop()

	results, err := pinger.GetLogContent()
	if err != nil {
		t.Fatal("Failed to get log results.", err)
	}

	if !strings.Contains(results, "Unable to determine site status - connect: network is unreachable") {
		t.Fatal("Failed to report Internet Access Error: ", results)
	}
}

// TestGetSites tests the database retrieval of the list of sites.
func TestGetSites(t *testing.T) {
	var sites database.Sites

	db, err := database.InitializeTestDB("")
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}
	defer db.Close()

	s1 := database.Site{Name: "Test", IsActive: true, URL: "http://www.test.com",
		PingIntervalSeconds: 60, TimeoutSeconds: 30}
	err = s1.CreateSite(db)
	if err != nil {
		t.Fatal("Failed to create new site:", err)
	}

	// Create the second site.
	s2 := database.Site{Name: "Test 2", IsActive: true, URL: "http://www.example.com",
		PingIntervalSeconds: 60, TimeoutSeconds: 30}
	err = s2.CreateSite(db)
	if err != nil {
		t.Fatal("Failed to create second site:", err)
	}

	sites, err = pinger.GetSites(db)
	// Verify the first site was Loaded with proper attributes.
	if !reflect.DeepEqual(s1.URL, sites[0].URL) || !reflect.DeepEqual(s1.IsActive, sites[0].IsActive) ||
		!reflect.DeepEqual(s1.Name, sites[0].Name) || !reflect.DeepEqual(s1.PingIntervalSeconds,
		sites[0].PingIntervalSeconds) || !reflect.DeepEqual(s1.TimeoutSeconds,
		sites[0].TimeoutSeconds) || !reflect.DeepEqual(s1.SiteID, sites[0].SiteID) {
		t.Fatal("First saved site not equal to input:\n", sites[0], s1)
	}

	// Verify the second site was Loaded with proper attributes.
	if !reflect.DeepEqual(s2.URL, sites[1].URL) || !reflect.DeepEqual(s2.IsActive, sites[1].IsActive) ||
		!reflect.DeepEqual(s2.Name, sites[1].Name) || !reflect.DeepEqual(s2.PingIntervalSeconds,
		sites[1].PingIntervalSeconds) || !reflect.DeepEqual(s2.TimeoutSeconds,
		sites[1].TimeoutSeconds) || !reflect.DeepEqual(s2.SiteID, sites[1].SiteID) {
		t.Fatal("First saved site not equal to input:\n", sites[1], s2)
	}
}

func TestCreatePingerLogError(t *testing.T) {
	var logFile = "/bogusFilePath/pinger.log"
	err := pinger.CreatePingerLog(logFile, true)
	if err == nil {
		t.Error("Creation of pinger log should throw error for bad path.")
	}
}
