package notifier

import (
	"log"
	"sync"

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
	var wg sync.WaitGroup
	log.Println("Sending Notification of Site Contacts about", n.Subject+"...")
	for _, c := range n.Site.Contacts {
		if c.SmsActive || c.EmailActive {
			// Notify contact
			wg.Add(1)
			go send(c, n.Message, n.Subject, &wg)
		} else {
			log.Println("No active contact methods for", c.Name)
		}
	}
	wg.Wait()
}

func send(c database.Contact, message string, subject string, wg *sync.WaitGroup) {
	log.Println("Sending notifications for", c.Name, subject, message)
	wg.Done()
}
