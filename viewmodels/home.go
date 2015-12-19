package viewmodels

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/turnkey-commerce/go-ping-sites/database"
)

// HomeViewModel holds the view information for the home.html template
type HomeViewModel struct {
	Error  error
	Title  string
	Active string
	Sites  []SiteViewModel
}

// SiteViewModel holds the required information about the site.
type SiteViewModel struct {
	Name     string
	Status   string
	HowLong  string
	CSSClass string
}

// GetHomeViewModel populates the items required by the home.html view
func GetHomeViewModel(db *sql.DB) HomeViewModel {
	result := HomeViewModel{
		Title:  "Go Ping Sites - Home",
		Active: "home",
		Error:  nil,
	}
	var sites database.Sites
	err := sites.GetActiveSites(db)
	if err != nil {
		result.Error = err
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
		siteVM.HowLong = fmt.Sprintf("%v", time.Since(site.LastStatusChange))

		result.Sites = append(result.Sites, *siteVM)
	}

	return result
}
