package viewmodels

import (
	"strconv"

	"github.com/apexskier/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/database"
)

// SitesEditViewModel holds the required information about the Sites to choose for editing.
type SitesEditViewModel struct {
	SiteID              int64  `valid:"-"`
	Name                string `valid:"ascii,required"`
	IsActive            bool   `valid:"-"`
	URL                 string `valid:"ascii,required"`
	PingIntervalSeconds string `valid:"int,required"`
	TimeoutSeconds      string `valid:"int,required"`
}

// SiteViewModel holds the view information for the site_edit.gohtml template
type SiteViewModel struct {
	Errors   map[string]string
	Title    string
	Site     SitesEditViewModel
	Contacts []database.Contact
	Nav      NavViewModel
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

// MapSiteVMtoDB maps the site view model properties to the site database properties.
func MapSiteVMtoDB(siteVM *SitesEditViewModel, site *database.Site) error {
	site.SiteID = siteVM.SiteID
	site.Name = siteVM.Name
	site.IsActive = siteVM.IsActive
	site.URL = siteVM.URL
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
	siteVM.PingIntervalSeconds = strconv.Itoa(site.PingIntervalSeconds)
	siteVM.TimeoutSeconds = strconv.Itoa(site.TimeoutSeconds)
}
