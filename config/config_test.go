package config_test

import (
	"testing"

	"github.com/turnkey-commerce/go-ping-sites/config"
)

func TestSmtpConfiguration(t *testing.T) {
	smtpSettings := config.Settings.SMTP

	if smtpSettings.EmailAddress != "yourusername@example.com" {
		t.Fatal("Config Email Address mismatch:\n", smtpSettings.EmailAddress)
	}

	if smtpSettings.Password != "yourpassword" {
		t.Fatal("Config Email Password mismatch:\n", smtpSettings.Password)
	}

	if smtpSettings.Server != "smtp.gmail.com" {
		t.Fatal("Config Email Server mismatch:\n", smtpSettings.Server)
	}

	if smtpSettings.Port != "587" {
		t.Fatal("Config Email Port mismatch:\n", smtpSettings.Port)
	}
}

func TestTwilioConfiguration(t *testing.T) {
	twilioSettings := config.Settings.Twilio

	if twilioSettings.AccountSid != "AccountSid" {
		t.Fatal("Config Twilio Account SID mismatch:\n", twilioSettings.AccountSid)
	}

	if twilioSettings.AuthToken != "AuthToken" {
		t.Fatal("Config Twilio Auth Token mismatch:\n", twilioSettings.AuthToken)
	}

	if twilioSettings.Number != "+15125551212" {
		t.Fatal("Config Twilio Number mismatch:\n", twilioSettings.Number)
	}
}
