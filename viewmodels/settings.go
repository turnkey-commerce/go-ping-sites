package viewmodels

import (
	"github.com/turnkey-commerce/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/database"
)

// SiteEditViewModel holds the required information about the site to choose for editing.
type SiteEditViewModel struct {
	SiteID              int64
	Name                string
	IsActive            string
	URL                 string
	PingIntervalSeconds int
	TimeoutSeconds      int
	NumContacts         int
}

// SettingsViewModel holds the view information for the settings.gohtml template
type SettingsViewModel struct {
	Error error
	Title string
	Sites []SiteEditViewModel
	Nav   NavViewModel
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
		if site.IsActive {
			siteVM.IsActive = "Y"
		} else {
			siteVM.IsActive = "N"
		}
		result.Sites = append(result.Sites, *siteVM)
	}

	return result
}
