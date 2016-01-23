package viewmodels

import (
	"github.com/apexskier/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/database"
)

// SitesEditViewModel holds the required information about the Sites to choose for editing.
type SitesEditViewModel struct {
	SiteID              int64  `valid:"-"`
	Name                string `valid:"ascii,required"`
	IsActive            bool   `valid:"-"`
	URL                 string `valid:"ascii,required"`
	PingIntervalSeconds int    `valid:"numeric,required"`
	TimeoutSeconds      int    `valid:"numeric,required"`
}

// SiteViewModel holds the view information for the site_edit.gohtml template
type SiteViewModel struct {
	Errors map[string]string
	Title  string
	Site   SitesEditViewModel
	Nav    NavViewModel
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
	siteVM.Name = site.Name
	siteVM.IsActive = site.IsActive
	siteVM.URL = site.URL
	siteVM.PingIntervalSeconds = site.PingIntervalSeconds
	siteVM.TimeoutSeconds = site.TimeoutSeconds

	result.Site = *siteVM

	return result
}
