package viewmodels

import (
	"fmt"

	"github.com/apexskier/httpauth"
	"github.com/dustin/go-humanize"
	"github.com/turnkey-commerce/go-ping-sites/database"
)

// HomeViewModel holds the view information for the home.gohtml template
type HomeViewModel struct {
	Title                      string
	Sites                      []SiteDashboardViewModel
	Nav                        NavViewModel
	Messages                   []string
	HasSiteWithNoStatusChanges bool
}

// SiteDashboardViewModel holds the required information about the site.
type SiteDashboardViewModel struct {
	SiteID             int64
	Name               string
	Status             string
	HowLong            string
	CSSClass           string
	LastChecked        string
	HasNoStatusChanges bool
}

// NavViewModel holds the information for the nav bar.
type NavViewModel struct {
	Active          string
	IsAuthenticated bool
	User            httpauth.UserData
	Messages        []string
}

// GetHomeViewModel populates the items required by the home.gohtml view
func GetHomeViewModel(sites database.Sites, isAuthenticated bool, user httpauth.UserData, messages []string) HomeViewModel {
	nav := NavViewModel{
		Active:          "home",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := HomeViewModel{
		Title:                      "Go Ping Sites - Home",
		Nav:                        nav,
		Messages:                   messages,
		HasSiteWithNoStatusChanges: false,
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
			siteVM.HasNoStatusChanges = true
			result.HasSiteWithNoStatusChanges = true
			// If LastStatusChange is zero then check if the first ping is also zero.
			if site.FirstPing.IsZero() {
				siteVM.HowLong = "Unknown"
			} else {
				// Use the first ping but set an asterisk to let the user know.
				siteVM.HowLong = fmt.Sprintf("%s", humanize.Time(site.FirstPing))
			}
		} else {
			siteVM.HasNoStatusChanges = false
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
