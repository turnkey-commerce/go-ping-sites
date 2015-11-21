package notifier_test

import (
	"reflect"
	"testing"

	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/notifier"
)

// TestNewNotifier tests building the pinger object.
func TestNewNotifier(t *testing.T) {
	site := getTestSite()
	n := notifier.NewNotifier(site, "Site 1 responding OK", "Site 1 Up", notifier.SendEmailMock, notifier.SendSmsMock)
	// Verify the first contact was Loaded with proper attributes and sorted last.
	if !reflect.DeepEqual(site.Contacts, n.Site.Contacts) {
		t.Fatal("Incoming site contacts are not the same as the notifier contacts:\n", site.Contacts, n.Site.Contacts)
	}
}

// TestNotify tests calling the Notifications.
func TestNotify(t *testing.T) {
	site := getTestSite()
	n := notifier.NewNotifier(site, "Site 1 responding OK", "Site 1 Up", notifier.SendEmailMock, notifier.SendSmsMock)
	n.Notify()
}

func getTestSite() database.Site {
	// Check that the contact got passed properly
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
