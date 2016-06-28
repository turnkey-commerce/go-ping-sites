package config_test

import (
	"testing"

	"github.com/turnkey-commerce/go-ping-sites/config"
)

func TestSmtpConfiguration(t *testing.T) {
	smtpSettings := config.Settings.SMTP

	if smtpSettings.EmailAddress != "yourusername@example.com" {
		t.Error("Config Email Address mismatch:\n", smtpSettings.EmailAddress)
	}

	if smtpSettings.Password != "yourpassword" {
		t.Error("Config Email Password mismatch:\n", smtpSettings.Password)
	}

	if smtpSettings.Server != "smtp.gmail.com" {
		t.Error("Config Email Server mismatch:\n", smtpSettings.Server)
	}

	if smtpSettings.Port != "587" {
		t.Fatal("Config Email Port mismatch:\n", smtpSettings.Port)
	}
}

func TestTwilioConfiguration(t *testing.T) {
	twilioSettings := config.Settings.Twilio

	if twilioSettings.AccountSid != "AccountSid" {
		t.Error("Config Twilio Account SID mismatch:\n", twilioSettings.AccountSid)
	}

	if twilioSettings.AuthToken != "AuthToken" {
		t.Error("Config Twilio Auth Token mismatch:\n", twilioSettings.AuthToken)
	}

	if twilioSettings.Number != "+15125551212" {
		t.Error("Config Twilio Number mismatch:\n", twilioSettings.Number)
	}
}

func TestWebsiteConfiguration(t *testing.T) {
	websiteSettings := config.Settings.Website

	if websiteSettings.HTTPPort != "8000" {
		t.Error("Config Website HTTPPort mismatch:\n", websiteSettings.HTTPPort)
	}

	if websiteSettings.CookieKey != "CookieEncryptionKey" {
		t.Error("Config Website HTTPPort mismatch:\n", websiteSettings.CookieKey)
	}

	if websiteSettings.SecureHTTPS != false {
		t.Error("Config Website SecureHTTPS mismatch:\n", websiteSettings.SecureHTTPS)
	}
}
