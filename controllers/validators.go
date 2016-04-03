package controllers

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

func validatePassword(allowMissingPassword bool, password string, password2 string, valErrors map[string]string) {
	if !allowMissingPassword && utf8.RuneCountInString(strings.TrimSpace(password)) < 6 {
		valErrors["Password"] = "Password must be at least 6 characters in length."
	} else if allowMissingPassword && utf8.RuneCountInString(strings.TrimSpace(password)) < 6 &&
		utf8.RuneCountInString(strings.TrimSpace(password)) > 0 {
		valErrors["Password"] = "Password must be at least 6 characters in length."
	}
	if password != password2 {
		valErrors["Password2"] = "Repeated Password must be the same as Password."
	}
}

func validateContact(contact *viewmodels.ContactsEditViewModel, valErrors map[string]string) {
	if contact.EmailActive && len(strings.TrimSpace(contact.EmailAddress)) == 0 {
		valErrors["EmailAddress"] = "Email Address must be provided if it is active."
	}
	if contact.SmsActive && len(strings.TrimSpace(contact.SmsNumber)) == 0 {
		valErrors["SmsNumber"] = "Text Message Number must be provided if it is active."
	}
	if contact.SmsActive && len(strings.TrimSpace(contact.SmsNumber)) > 0 {
		var validSmsNumber = regexp.MustCompile(`^\+?[1-9][0-9]{1,14}$`)
		if !validSmsNumber.MatchString(contact.SmsNumber) {
			valErrors["SmsNumber"] = "The Text Message Number must be provided in E.164 format. For example in the USA it would be +15125551212."
		}
	}
}
