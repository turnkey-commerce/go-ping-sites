package viewmodels

import (
	"fmt"

	"github.com/apexskier/httpauth"
	"github.com/dustin/go-humanize"
	"github.com/turnkey-commerce/go-ping-sites/database"
)

// HomeViewModel holds the view information for the home.gohtml template
type HomeViewModel struct {
	Error    error
	Title    string
	Sites    []SiteDashboardViewModel
	Nav      NavViewModel
	Messages []string
}

// SiteDashboardViewModel holds the required information about the site.
type SiteDashboardViewModel struct {
	SiteID      int64
	Name        string
	Status      string
	HowLong     string
	CSSClass    string
	LastChecked string
}

// NavViewModel holds the information for the nav bar.
type NavViewModel struct {
	Active          string
	IsAuthenticated bool
	User            httpauth.UserData
	Messages        []string
}

// GetHomeViewModel populates the items required by the home.gohtml view
func GetHomeViewModel(sites database.Sites, isAuthenticated bool, user httpauth.UserData, messages []string, err error) HomeViewModel {
	nav := NavViewModel{
		Active:          "home",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := HomeViewModel{
		Title:    "Go Ping Sites - Home",
		Error:    err,
		Nav:      nav,
		Messages: messages,
	}

	for _, site := range sites {
		siteVM := new(SiteDashboardViewModel)
		siteVM.Name = site.Name
		siteVM.SiteID = site.SiteID

		if site.IsSiteUp {
			siteVM.Status = "Up"
			siteVM.CSSClass = "success"
		} else {
			siteVM.Status = "Down"
			siteVM.CSSClass = "danger"
		}

		if site.LastStatusChange.IsZero() {
			siteVM.HowLong = "Unknown"
		} else {
			siteVM.HowLong = fmt.Sprintf("%s", humanize.Time(site.LastStatusChange))
		}

		if site.LastPing.IsZero() {
			siteVM.LastChecked = "Never"
		} else {
			siteVM.LastChecked = fmt.Sprintf("%s", humanize.Time(site.LastPing))
		}

		result.Sites = append(result.Sites, *siteVM)
	}

	return result
}
