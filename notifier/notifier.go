package notifier

import (
	"log"
	"net/smtp"
	"sync"

	"github.com/sfreiberg/gotwilio"

	"github.com/turnkey-commerce/go-ping-sites/config"
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
	var err error
	log.Println("Sending notifications for", c.Name, subject, message)
	if c.EmailActive && len(c.EmailAddress) > 0 {
		err = sendEmail(c.EmailAddress, message, subject)
		if err != nil {
			log.Println("Error sending email:", err)
		}
	}

	if c.SmsActive && len(c.SmsNumber) > 0 {
		err = sendSms(c.SmsNumber, message)
		if err != nil {
			log.Println("Error sending SMS:", err)
		}
	}

	wg.Done()
}

func sendEmail(recipient string, message string, subject string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", config.Settings.SMTP.EmailAddress, config.Settings.SMTP.Password,
		config.Settings.SMTP.Server)
	server := config.Settings.SMTP.Server + ":" + config.Settings.SMTP.Port
	to := []string{recipient}
	from := "sender@example.org"
	msg := []byte("Subject: " + subject + "\r\n\r\n" +
		message + "\r\n")
	err := smtp.SendMail(server, auth, from, to, msg)
	if err != nil {
		return err
	}
	return nil
}

func sendSms(smsNumber string, message string) error {
	twilio := gotwilio.NewTwilioClient(config.Settings.Twilio.AccountSid, config.Settings.Twilio.AuthToken)
	from := config.Settings.Twilio.Number

	_, _, err := twilio.SendSMS(from, smsNumber, message, "", "")
	if err != nil {
		return err
	}
	return nil
}
