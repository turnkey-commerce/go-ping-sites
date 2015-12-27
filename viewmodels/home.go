package viewmodels

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/turnkey-commerce/go-ping-sites/database"
)

// HomeViewModel holds the view information for the home.gohtml template
type HomeViewModel struct {
	Error  error
	Title  string
	Active string
	Sites  []SiteViewModel
}

// SiteViewModel holds the required information about the site.
type SiteViewModel struct {
	Name        string
	Status      string
	HowLong     string
	CSSClass    string
	LastChecked string
}

// GetHomeViewModel populates the items required by the home.gohtml view
func GetHomeViewModel(sites database.Sites, err error) HomeViewModel {
	result := HomeViewModel{
		Title:  "Go Ping Sites - Home",
		Active: "home",
		Error:  err,
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
