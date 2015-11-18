package notifier

import (
	"log"

	"github.com/turnkey-commerce/go-ping-sites/database"
)

// Notifier sends the notifications to the recipients on a status change.
type Notifier struct {
	Site    database.Site
	Message string
	Subject string
}

// NewNotifier returns a new Notifier object to perform notifications about status change
func NewNotifier(site database.Site, message string, subject string) *Notifier {
	n := Notifier{Site: site, Message: message, Subject: subject}
	return &n
}

// Notify starts the notification for each contact for the site.
func (n *Notifier) Notify() {
	log.Println("Sending Notification of Site Contacts about", n.Subject+"...")
	for _, c := range n.Site.Contacts {
		if c.SmsActive || c.EmailActive {
			// Notify contacts
			log.Println("Sending notifications for", c.Name)
		} else {
			log.Println("No active contact methods for", c.Name)
		}
	}
}
