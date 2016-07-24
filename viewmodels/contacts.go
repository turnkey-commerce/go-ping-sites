package viewmodels

import (
	"html/template"

	"github.com/apexskier/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/database"
)

// ContactsEditViewModel holds the required information about the Contacts to choose for editing.
type ContactsEditViewModel struct {
	ContactID     int64   `valid:"-"`
	Name          string  `valid:"ascii,required"`
	EmailAddress  string  `valid:"email"`
	SmsNumber     string  `valid:"-"`
	SmsActive     bool    `valid:"-"`
	EmailActive   bool    `valid:"-"`
	SelectedSites []int64 `valid:"-"`
	SiteCount     int     `valid:"-"`
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
	Errors    map[string]string
	Title     string
	Contact   ContactsEditViewModel
	AllSites  []ContactsAllSitesViewModel
	Nav       NavViewModel
	CsrfField template.HTML
}

// ContactsAllSitesViewModel has all of the sites available to assign the contact.
type ContactsAllSitesViewModel struct {
	SiteID     int64
	Name       string
	IsActive   bool
	URL        string
	IsAssigned bool
}

// GetContactsViewModel populates the items required by the contacts.gohtml view
func GetContactsViewModel(contacts database.Contacts, isAuthenticated bool,
	user httpauth.UserData, err error) ContactsViewModel {
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
		contactVM.SiteCount = contact.SiteCount

		result.Contacts = append(result.Contacts, *contactVM)
	}

	return result
}

// EditContactViewModel populates the items required by the user_contact.gohtml view
func EditContactViewModel(formContact *ContactsEditViewModel, allSites database.Sites,
	isAuthenticated bool, user httpauth.UserData, errors map[string]string) ContactViewModel {
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
	result.AllSites = PopulateAllSitesVM(allSites, formContact.SelectedSites,
		false)

	return result
}

// NewContactViewModel populates the items required by the user_contact.gohtml view
func NewContactViewModel(formContact *ContactsEditViewModel, allSites database.Sites,
	selectAllSites bool, isAuthenticated bool, user httpauth.UserData,
	errors map[string]string) ContactViewModel {
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
	result.AllSites = PopulateAllSitesVM(allSites, formContact.SelectedSites,
		selectAllSites)

	return result
}

// PopulateAllSitesVM returns the view model for all of the sites
func PopulateAllSitesVM(allSites database.Sites, selectedSiteIDs []int64,
	selectAllSites bool) []ContactsAllSitesViewModel {
	var allSitesVM = []ContactsAllSitesViewModel{}
	for _, site := range allSites {
		hasMatch := false
		if selectAllSites {
			hasMatch = true
		} else {
			for _, siteID := range selectedSiteIDs {
				if siteID == site.SiteID {
					hasMatch = true
					break
				}
			}
		}
		siteVM := ContactsAllSitesViewModel{
			SiteID:     site.SiteID,
			Name:       site.Name,
			IsActive:   site.IsActive,
			URL:        site.URL,
			IsAssigned: hasMatch,
		}
		allSitesVM = append(allSitesVM, siteVM)
	}
	return allSitesVM
}
