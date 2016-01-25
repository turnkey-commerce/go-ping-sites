package viewmodels

import (
	"strconv"

	"github.com/apexskier/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/database"
)

// SitesEditViewModel holds the required information about the Sites to choose for editing.
// The PingIntervalSeconds and TimeoutSeconds are strings to allow the form validation.
type SitesEditViewModel struct {
	SiteID              int64  `valid:"-"`
	Name                string `valid:"ascii,required"`
	IsActive            bool   `valid:"-"`
	URL                 string `valid:"ascii,required"`
	PingIntervalSeconds string `valid:"int,required"`
	TimeoutSeconds      string `valid:"int,required"`
}

// SitesAllContactsViewModel has all of the sites available and carries whether
// the contact is part of the Site itself via the IsAssigned property
type SitesAllContactsViewModel struct {
	ContactID    int64
	IsAssigned   bool
	Name         string
	EmailAddress string
	SmsNumber    string
	SmsActive    bool
	EmailActive  bool
}

// SiteContactsSelectedViewModel holds the selections when contacts are changed.
// The existing SiteContacts are also containd in SiteContats
type SiteContactsSelectedViewModel struct {
	SiteID           int64
	SelectedContacts []int64
	SiteContacts     []int64
}

// SiteViewModel holds the view information for the site_edit.gohtml template
type SiteViewModel struct {
	Errors      map[string]string
	Title       string
	Site        SitesEditViewModel
	Contacts    []database.Contact
	AllContacts []SitesAllContactsViewModel
	Nav         NavViewModel
}

// GetSiteDetailsViewModel populates the items required by the site_details.gohtml view
func GetSiteDetailsViewModel(site *database.Site, isAuthenticated bool,
	user httpauth.UserData) SiteViewModel {
	nav := NavViewModel{
		Active:          "settings",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := SiteViewModel{
		Title: "Go Ping Sites - Settings - Site Details",
		Nav:   nav,
	}

	siteVM := new(SitesEditViewModel)
	MapSiteDBtoVM(site, siteVM)
	result.Site = *siteVM
	result.Contacts = site.Contacts

	return result
}

// EditSiteViewModel populates the items required by the site_edit.gohtml view
func EditSiteViewModel(siteVM *SitesEditViewModel, isAuthenticated bool,
	user httpauth.UserData, errors map[string]string) SiteViewModel {
	nav := NavViewModel{
		Active:          "settings",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := SiteViewModel{
		Title:  "Go Ping Sites - Settings - Edit Site",
		Nav:    nav,
		Errors: errors,
	}

	result.Site = *siteVM
	return result
}

// NewSiteViewModel populates the items required by the site_new.gohtml view
func NewSiteViewModel(siteVM *SitesEditViewModel, isAuthenticated bool,
	user httpauth.UserData, errors map[string]string) SiteViewModel {
	nav := NavViewModel{
		Active:          "settings",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := SiteViewModel{
		Title:  "Go Ping Sites - Settings - New Site",
		Nav:    nav,
		Errors: errors,
	}
	result.Site = *siteVM
	return result
}

// SiteChangeContactsViewModel populates the items required by the
// site_change_contacts.gohtml template
func SiteChangeContactsViewModel(site *database.Site, allContacts database.Contacts, isAuthenticated bool,
	user httpauth.UserData) SiteViewModel {
	nav := NavViewModel{
		Active:          "settings",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := SiteViewModel{
		Title: "Go Ping Sites - Settings - Change Contacts",
		Nav:   nav,
	}

	siteVM := new(SitesEditViewModel)
	MapSiteDBtoVM(site, siteVM)
	result.Site = *siteVM
	result.Contacts = site.Contacts
	result.AllContacts = PopulateAllContactsVM(allContacts, site.Contacts)

	return result
}

// MapSiteVMtoDB maps the site view model properties to the site database properties.
func MapSiteVMtoDB(siteVM *SitesEditViewModel, site *database.Site) error {
	site.SiteID = siteVM.SiteID
	site.Name = siteVM.Name
	site.IsActive = siteVM.IsActive
	site.URL = siteVM.URL
	// Conversion on these two is necessary because they are a string in the
	// view model to allow the validation to work
	pingInterval, err := strconv.Atoi(siteVM.PingIntervalSeconds)
	if err != nil {
		return err
	}
	site.PingIntervalSeconds = pingInterval
	timeout, err := strconv.Atoi(siteVM.TimeoutSeconds)
	if err != nil {
		return err
	}
	site.TimeoutSeconds = timeout

	return nil
}

// MapSiteDBtoVM maps the site database properties to the site view model properties.
func MapSiteDBtoVM(site *database.Site, siteVM *SitesEditViewModel) {
	siteVM.SiteID = site.SiteID
	siteVM.Name = site.Name
	siteVM.IsActive = site.IsActive
	siteVM.URL = site.URL
	// Conversion on these two is necessary because they are a string in the
	// view model to allow the validation to work
	siteVM.PingIntervalSeconds = strconv.Itoa(site.PingIntervalSeconds)
	siteVM.TimeoutSeconds = strconv.Itoa(site.TimeoutSeconds)
}

// PopulateAllContactsVM returns the view model for the contacts with the ones
// assigned to the site having IsAssigned set to true.
func PopulateAllContactsVM(allContacts database.Contacts,
	siteContacts database.Contacts) []SitesAllContactsViewModel {
	var allContactsVM = []SitesAllContactsViewModel{}
	for _, contact := range allContacts {
		hasMatch := false
		for _, siteContact := range siteContacts {
			if siteContact.ContactID == contact.ContactID {
				hasMatch = true
				break
			}
		}
		contactVM := SitesAllContactsViewModel{
			ContactID:    contact.ContactID,
			Name:         contact.Name,
			IsAssigned:   hasMatch,
			EmailAddress: contact.EmailAddress,
			EmailActive:  contact.EmailActive,
			SmsNumber:    contact.SmsNumber,
			SmsActive:    contact.SmsActive,
		}
		allContactsVM = append(allContactsVM, contactVM)
	}
	return allContactsVM
}
