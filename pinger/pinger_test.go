package pinger_test

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"strings"
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

	results, err := getLogContent()
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

// TestStartPinger starts up the pinger and then stops it after 10 seconds
func TestStartPinger(t *testing.T) {
	p := pinger.NewPinger(nil, getSites, requestURL)
	p.Start()
	time.Sleep(10 * time.Second)
	p.Stop()

	results, err := getLogContent()
	if err != nil {
		t.Fatal("Failed to get log results.", err)
	}

	if !strings.Contains(results, "Client.Timeout") {
		t.Fatal("Failed to report timeout error.")
	}
	if !strings.Contains(results, "Test 3 Paused") {
		t.Fatal("Failed to report paused site.")
	}
	if !strings.Contains(results, "Error - HTTP Status Code") {
		t.Fatal("Failed to report bad HTTP Status Code.")
	}
}

// TestStartPinger starts up the pinger and then stops it after 10 seconds
func TestStartEmptySitesPinger(t *testing.T) {
	p := pinger.NewPinger(nil, getEmptySites, requestURL)
	p.Start()
	time.Sleep(1 * time.Second)
	p.Stop()

	results, err := getLogContent()
	if err != nil {
		t.Fatal("Failed to get log results.", err)
	}

	if !strings.Contains(results, "No active sites set up for pinging.") {
		t.Fatal("Failed to report empty sites.")
	}
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
		PingIntervalSeconds: 2, TimeoutSeconds: 1}
	// Create the second site.
	s2 := database.Site{Name: "Test 2", IsActive: true, URL: "http://www.github.com",
		PingIntervalSeconds: 5, TimeoutSeconds: 2}
	// Create the third site as not active.
	s3 := database.Site{Name: "Test 3", IsActive: false, URL: "http://www.test.com",
		PingIntervalSeconds: 5, TimeoutSeconds: 2}
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

func getEmptySites(db *sql.DB) (database.Sites, error) {
	var sites database.Sites
	return sites, nil
}

// getLogContent reads the results of the log file for verification.
func getLogContent() (string, error) {
	dat, err := ioutil.ReadFile("pinger.log")
	if err != nil {
		return "", err
	}
	results := string(dat)
	return results, nil
}
