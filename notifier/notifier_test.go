package notifier_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/notifier"
	"github.com/turnkey-commerce/go-ping-sites/pinger"
)

// TestNewNotifier tests building the pinger object.
func TestNewNotifier(t *testing.T) {
	site := getTestSite()
	n := notifier.NewNotifier(site, "Site 1 responding OK", "Site 1 Up", notifier.SendEmailMock, notifier.SendSmsMock)
	// Verify the first contact was Loaded with proper attributes and sorted last.
	if !reflect.DeepEqual(site.Contacts, n.Site.Contacts) {
		t.Error("Incoming site contacts are not the same as the notifier contacts:\n", site.Contacts, n.Site.Contacts)
	}
}

// TestNotify tests calling the Notifications successfully.
func TestNotify(t *testing.T) {
	pinger.CreatePingerLog("", true)
	site := getTestSite()
	n := notifier.NewNotifier(site, "Site 1 responding OK", "Site 1 Up", notifier.SendEmailMock, notifier.SendSmsMock)
	n.Notify()

	results, err := pinger.GetLogContent()
	if err != nil {
		t.Fatal("Failed to get log results.", err)
	}

	if !strings.Contains(results, "Sending notifications for Jack Contact Site 1 Up Site 1 responding OK") {
		t.Error("Failed to report successful send to Jack Contact.")
	}

	if !strings.Contains(results, "Sending notifications for Joe Contact Site 1 Up Site 1 responding OK") {
		t.Error("Failed to report successful send to Joe Contact.")
	}
}

// TestNotifyError tests calling the Notifications with errors on each send method.
func TestNotifyError(t *testing.T) {
	site := getTestSite()
	n := notifier.NewNotifier(site, "Site 1 responding OK", "Site 1 Up", notifier.SendEmailErrorMock, notifier.SendSmsErrorMock)
	n.Notify()

	results, err := pinger.GetLogContent()
	if err != nil {
		t.Fatal("Failed to get log results.", err)
	}

	if !strings.Contains(results, "Error sending SMS: Error - no response from server.") {
		t.Error("Failed to report error in SMS send to Jack Contact.")
	}

	if !strings.Contains(results, "Error sending email: Error - no response from server.") {
		t.Error("Failed to report error in email send to Joe Contact.")
	}
}

// Create the struct for the Site and its contacts used for testing.
func getTestSite() database.Site {
	s1 := database.Site{Name: "Test", IsActive: true, URL: "http://www.google.com",
		PingIntervalSeconds: 2, TimeoutSeconds: 1}

	// Create first contact
	c1 := database.Contact{Name: "Joe Contact", EmailAddress: "joe@test.com", SmsNumber: "5125551212",
		SmsActive: false, EmailActive: true}
	// Create second contact
	c2 := database.Contact{Name: "Jack Contact", EmailAddress: "jack@test.com", SmsNumber: "5125551213",
		SmsActive: true, EmailActive: false}

	// Add the contacts to the sites
	s1.Contacts = append(s1.Contacts, c1, c2)

	return s1
}
