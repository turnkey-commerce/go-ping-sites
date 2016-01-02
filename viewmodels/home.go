package viewmodels

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/turnkey-commerce/go-ping-sites/database"
)

// HomeViewModel holds the view information for the home.gohtml template
type HomeViewModel struct {
	Error           error
	Title           string
	Active          string
	IsAuthenticated bool
	Sites           []SiteViewModel
	Nav             NavViewModel
	Messages        []string
}

// SiteViewModel holds the required information about the site.
type SiteViewModel struct {
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
	Messages        []string
}

// GetHomeViewModel populates the items required by the home.gohtml view
func GetHomeViewModel(sites database.Sites, isAuthenticated bool, messages []string, err error) HomeViewModel {
	nav := NavViewModel{
		Active:          "home",
		IsAuthenticated: isAuthenticated,
	}

	result := HomeViewModel{
		Title:    "Go Ping Sites - Home",
		Active:   "home",
		Error:    err,
		Nav:      nav,
		Messages: messages,
	}

	for _, site := range sites {
		siteVM := new(SiteViewModel)
		siteVM.Name = site.Name

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
