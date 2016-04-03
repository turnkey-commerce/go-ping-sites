package controllers

import (
	"strings"
	"testing"

	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

func TestValidatePasswordValid(t *testing.T) {
	u := new(viewmodels.UsersEditViewModel)
	u.Email = "jack@example.com"
	u.Username = ""
	valErrors := validateUserForm(u, true)
	if !strings.Contains(valErrors["Username"], "non zero value required") {
		t.Error("UserName should show error for required.")
	}

	u.Username = "jack"
	u.Email = "jack.example.com"
	valErrors = validateUserForm(u, true)
	if !strings.Contains(valErrors["Email"], "does not validate as email") {
		t.Error("Email Address Validation should show error for invalid email.")
	}

	u.Username = "jack"
	u.Email = "jack@example.com"
	u.Password = "short"
	valErrors = validateUserForm(u, false)
	if !strings.Contains(valErrors["Password"], "Password must be at least 6 characters in length") {
		t.Error("Password Validation should show error for short password.")
	}

	u.Password = "short"
	valErrors = validateUserForm(u, false)
	if !strings.Contains(valErrors["Password"], "Password must be at least 6 characters in length") {
		t.Error("Password Validation should show error for short password.")
	}

	u.Username = "jack"
	u.Email = "jack@example.com"
	u.Password = "long enough"
	u.Password2 = "long enough2"
	valErrors = validateUserForm(u, false)
	if !strings.Contains(valErrors["Password2"], "Repeated Password must be the same as Password.") {
		t.Error("Password Validation should show error for short password.")
	}

	u.Username = "jack"
	u.Email = "jack@example.com"
	u.Password = "long enough"
	u.Password2 = "long enough"
	valErrors = validateUserForm(u, false)
	if len(valErrors) > 0 {
		t.Error("No errors should be flagged for the set of inputs.")
	}
}

func TestValidateContactValid(t *testing.T) {
	c := new(viewmodels.ContactsEditViewModel)
	c.Name = ""
	c.SmsActive = false
	c.EmailActive = false
	c.SmsNumber = ""
	c.EmailAddress = ""
	valErrors := validateContactForm(c)
	if !strings.Contains(valErrors["Name"], "non zero value required") {
		t.Error("Name should show error for required.")
	}

	c.Name = "Jane Doe"
	c.SmsActive = true
	c.EmailActive = false
	c.SmsNumber = ""
	valErrors = validateContactForm(c)
	if !strings.Contains(valErrors["SmsNumber"], "Text Message Number must be provided") {
		t.Error("Text Message Validation should show error for required if active.")
	}

	c.SmsActive = false
	c.EmailActive = true
	c.EmailAddress = ""
	valErrors = validateContactForm(c)
	if !strings.Contains(valErrors["EmailAddress"], "Email Address must be provided") {
		t.Error("Email Address Validation should show error for required if active.")
	}

	c.SmsActive = false
	c.EmailActive = true
	c.EmailAddress = "jack.example.com"
	valErrors = validateContactForm(c)
	if !strings.Contains(valErrors["EmailAddress"], "does not validate as email") {
		t.Error("Email Address Validation should show error for invalid email.")
	}

	c.SmsActive = true
	c.EmailActive = false
	c.SmsNumber = "Foobar"
	valErrors = validateContactForm(c)
	if !strings.Contains(valErrors["SmsNumber"], "The Text Message Number must be provided in E.164 format") {
		t.Error("Text Message Validation should show error for incorrect number.")
	}

	c.SmsNumber = "1234567890123456"
	valErrors = validateContactForm(c)
	if !strings.Contains(valErrors["SmsNumber"], "The Text Message Number must be provided in E.164 format") {
		t.Error("Text Message Validation should show error for incorrect number.")
	}

	c.SmsNumber = "+15127712936"
	c.Name = "Jane Doe"
	c.SmsActive = true
	c.EmailActive = true
	c.EmailAddress = "jane@example.com"
	valErrors = validateContactForm(c)
	if len(valErrors) > 0 {
		t.Error("No errors should be flagged for the set of inputs.")
	}
}
