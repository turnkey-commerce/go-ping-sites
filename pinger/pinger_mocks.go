package pinger

import (
	"database/sql"
	"errors"
	"time"

	"github.com/turnkey-commerce/go-ping-sites/database"
)

// hitCount is used to vary the outcome of the mock RequestURL
var hitCount int

// ResetHitCount sets the hitcount back to 0 for the tests.
func ResetHitCount() {
	hitCount = 0
}

// RequestURLMock is a mock of the URL request that pings the site.
func RequestURLMock(url string, timeout int) (string, int, time.Duration, error) {
	var responseTime = 300 * time.Millisecond
	hitCount++
	// The hitCount allows to vary the response of the request.
	if url == "http://www.github.com" && hitCount < 4 {
		return "", 0, responseTime, errors.New("(Client.Timeout exceeded while awaiting headers)")
	} else if url == "http://www.github.com" {
		return "Hello", 200, responseTime, nil
	}
	return "Hello", 300, responseTime, nil
}

// RequestURLContentMock is a mock of the URL requests for checking content.
func RequestURLContentMock(url string, timeout int) (string, int, time.Duration, error) {
	var responseTime = 300 * time.Millisecond
	// The hitCount allows to vary the response of the request.
	if url == "http://www.github.com" {
		return "Bad response text", 200, responseTime, nil
	} else if url == "http://www.google.com" {
		return "Good response text", 200, responseTime, nil
	}
	return "", 200, responseTime, nil
}

// RequestURLBadInternetAccessMock mocks the condition where the outgoing Internet connection is down.
func RequestURLBadInternetAccessMock(url string, timeout int) (string, int, time.Duration, error) {
	var responseTime = 300 * time.Millisecond
	return "", 0, responseTime, InternetAccessError{msg: "connect: network is unreachable"}
}

// GetSitesMock is a mock of the SQL query to get the sites for pinging
func GetSitesMock(db *sql.DB) (database.Sites, error) {
	var sites database.Sites
	// Create the first site.
	s1 := database.Site{Name: "Test", IsActive: true, URL: "http://www.google.com",
		PingIntervalSeconds: 1, TimeoutSeconds: 1, IsSiteUp: true}
	// Create the second site.
	s2 := database.Site{Name: "Test 2", IsActive: true, URL: "http://www.github.com",
		PingIntervalSeconds: 2, TimeoutSeconds: 2, IsSiteUp: true}
	// Create the third site as not active.
	s3 := database.Site{Name: "Test 3", IsActive: false, URL: "http://www.test.com",
		PingIntervalSeconds: 2, TimeoutSeconds: 2}
	// Contacts are deliberately set as false for SmsActive and EmailActive so as not to trigger Notifier
	c1 := database.Contact{Name: "Joe Contact", EmailAddress: "joe@test.com", SmsNumber: "5125551212",
		SmsActive: false, EmailActive: false}
	c2 := database.Contact{Name: "Jack Contact", EmailAddress: "jack@test.com", SmsNumber: "5125551213",
		SmsActive: false, EmailActive: false}
	// Add the contacts to the sites
	s1.Contacts = append(s1.Contacts, c1, c2)
	s2.Contacts = append(s2.Contacts, c1)
	s3.Contacts = append(s3.Contacts, c1)

	sites = append(sites, s1, s2, s3)
	return sites, nil
}

// GetSitesContentMock is a mock of the SQL query to get the sites for testing content checks.
func GetSitesContentMock(db *sql.DB) (database.Sites, error) {
	var sites database.Sites
	// Create the first site.
	s1 := database.Site{Name: "Test", IsActive: true, URL: "http://www.google.com",
		PingIntervalSeconds: 1, TimeoutSeconds: 1, IsSiteUp: true,
		ContentExpected: "Good response text", ContentUnexpected: "Bad response text"}
	// Create the second site.
	s2 := database.Site{Name: "Test 2", IsActive: true, URL: "http://www.github.com",
		PingIntervalSeconds: 1, TimeoutSeconds: 1, IsSiteUp: true,
		ContentExpected: "", ContentUnexpected: "Bad response text"}
	// Create the third site.
	s3 := database.Site{Name: "Test 3", IsActive: true, URL: "http://www.example.com",
		PingIntervalSeconds: 1, TimeoutSeconds: 1, IsSiteUp: true,
		ContentExpected: "Good response text", ContentUnexpected: "Bad response text"}

	sites = append(sites, s1, s2, s3)
	return sites, nil
}

// GetEmptySitesMock is a mock of the SQL query to get the sites for pinging
// In this case the method returns an empty list of sites.
func GetEmptySitesMock(db *sql.DB) (database.Sites, error) {
	var sites database.Sites
	return sites, nil
}

// GetSitesErrorMock is a mock of the SQL query to get the sites for pinging
// In this case it returns an error when getting the sites.
func GetSitesErrorMock(db *sql.DB) (database.Sites, error) {
	return nil, errors.New("Timeout accessing the SQL database.")
}
