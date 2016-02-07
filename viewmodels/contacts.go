package viewmodels

import (
	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/apexskier/httpauth"
)

// ContactsEditViewModel holds the required information about the Contacts to choose for editing.
type ContactsEditViewModel struct {
	ContactID    int64  `valid:"-"`
	Name         string `valid:"ascii,required"`
	EmailAddress string `valid:"email,required"`
	SmsNumber    string `valid:"alphanum,required"`
	SmsActive    bool   `valid:"-"`
	EmailActive  bool   `valid:"-"`
}

// ContactsViewModel holds the view information for the contacts.gohtml template
type ContactsViewModel struct {
	Error    error
	Title    string
	Contacts []ContactsEditViewModel
	Nav      NavViewModel
}

// ContactViewModel holds the view information for the contact_edit.gohtml template
type ContactViewModel struct {
	Errors  map[string]string
	Title   string
	Contact ContactsEditViewModel
	Nav     NavViewModel
}

// GetContactsViewModel populates the items required by the contacts.gohtml view
func GetContactsViewModel(contacts database.Contacts, isAuthenticated bool, user httpauth.UserData, err error) ContactsViewModel {
	nav := NavViewModel{
		Active:          "settings",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := ContactsViewModel{
		Title: "Go Ping Sites - Settings - Contacts",
		Nav:   nav,
	}

	for _, contact := range contacts {
		contactVM := new(ContactsEditViewModel)
		contactVM.ContactID = contact.ContactID
		contactVM.Name = contact.Name
		contactVM.EmailAddress = contact.EmailAddress
		contactVM.EmailActive = contact.EmailActive
		contactVM.SmsNumber = contact.SmsNumber
		contactVM.SmsActive = contact.SmsActive

		result.Contacts = append(result.Contacts, *contactVM)
	}

	return result
}

// EditContactViewModel populates the items required by the user_contact.gohtml view
func EditContactViewModel(formContact *ContactsEditViewModel, isAuthenticated bool,
	user httpauth.UserData, errors map[string]string) ContactViewModel {
	nav := NavViewModel{
		Active:          "settings",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := ContactViewModel{
		Title:  "Go Ping Sites - Settings - Edit Contact",
		Nav:    nav,
		Errors: errors,
	}

	contactVM := new(ContactsEditViewModel)
	contactVM.ContactID = formContact.ContactID
	contactVM.Name = formContact.Name
	contactVM.EmailAddress = formContact.EmailAddress
	contactVM.EmailActive = formContact.EmailActive
	contactVM.SmsNumber = formContact.SmsNumber
	contactVM.SmsActive = formContact.SmsActive

	result.Contact = *contactVM

	return result
}

// NewContactViewModel populates the items required by the user_contact.gohtml view
func NewContactViewModel(formContact *ContactsEditViewModel, isAuthenticated bool,
	user httpauth.UserData, errors map[string]string) ContactViewModel {
	nav := NavViewModel{
		Active:          "settings",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := ContactViewModel{
		Title:  "Go Ping Sites - Settings - New Contact",
		Nav:    nav,
		Errors: errors,
	}

	contactVM := new(ContactsEditViewModel)
	contactVM.ContactID = formContact.ContactID
	contactVM.Name = formContact.Name
	contactVM.EmailAddress = formContact.EmailAddress
	contactVM.EmailActive = formContact.EmailActive
	contactVM.SmsNumber = formContact.SmsNumber
	contactVM.SmsActive = formContact.SmsActive

	result.Contact = *contactVM

	return result
}
