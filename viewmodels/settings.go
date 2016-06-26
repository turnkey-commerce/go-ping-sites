package viewmodels

import (
	"github.com/apexskier/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/database"
	"html/template"
)

// SiteEditViewModel holds the required information about the site to choose for editing.
type SiteEditViewModel struct {
	SiteID              int64
	Name                string
	IsActive            bool
	URL                 string
	PingIntervalSeconds int
	TimeoutSeconds      int
	NumContacts         int
}

// SettingsViewModel holds the view information for the settings.gohtml template
type SettingsViewModel struct {
	Error     error
	Title     string
	Sites     []SiteEditViewModel
	Nav       NavViewModel
	CsrfField template.HTML
}

// GetSettingsViewModel populates the items required by the settings.gohtml view
func GetSettingsViewModel(sites database.Sites, isAuthenticated bool, user httpauth.UserData, err error) SettingsViewModel {
	nav := NavViewModel{
		Active:          "settings",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := SettingsViewModel{
		Title: "Go Ping Sites - Settings",
		Nav:   nav,
	}

	for _, site := range sites {
		siteVM := new(SiteEditViewModel)
		siteVM.Name = site.Name
		siteVM.URL = site.URL
		siteVM.SiteID = site.SiteID
		siteVM.PingIntervalSeconds = site.PingIntervalSeconds
		siteVM.TimeoutSeconds = site.TimeoutSeconds
		siteVM.NumContacts = len(site.Contacts)
		siteVM.IsActive = site.IsActive
		result.Sites = append(result.Sites, *siteVM)
	}

	return result
}
