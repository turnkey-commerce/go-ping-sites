package viewmodels_test

import (
	"testing"
	"time"

	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
	"github.com/turnkey-commerce/httpauth"
)

// TestGetHomeViewModel tests the view model is created as expected.
func TestGetHomeViewModel(t *testing.T) {
	sites := database.Sites{}
	user := httpauth.UserData{}

	// First site has no last status change.
	sites = append(sites, database.Site{Name: "Test 1", IsSiteUp: true})
	// Second site was down 2 hours ago.
	now := time.Now()
	twoHoursAgo := now.Add(-2 * time.Hour)
	sites = append(sites, database.Site{Name: "Test 2", IsSiteUp: false, LastStatusChange: twoHoursAgo})
	// Second site was up 2 days ago.
	twodaysAgo := now.Add(-48 * time.Hour)
	sites = append(sites, database.Site{Name: "Test 3", IsSiteUp: true, LastStatusChange: twodaysAgo})

	result := viewmodels.GetHomeViewModel(sites, false, user, nil, nil)
	if result.Nav.Active != "home" {
		t.Error("Home View Model Active returned incorrect value")
	}

	if result.Error != nil {
		t.Error("Home View Model Error should be nil")
	}

	if result.Title != "Go Ping Sites - Home" {
		t.Error("Home View Model Sites returned incorrect value")
	}

	if result.Sites[0].CSSClass != "success" || result.Sites[0].Status != "Up" ||
		result.Sites[0].HowLong != "Unknown" || result.Sites[0].Name != "Test 1" {
		t.Error("First site returned incorrect values")
	}

	if result.Sites[1].CSSClass != "danger" || result.Sites[1].Status != "Down" ||
		result.Sites[1].HowLong != "2 hours ago" || result.Sites[1].Name != "Test 2" {
		t.Error("Second site returned incorrect values")
	}

	if result.Sites[2].CSSClass != "success" || result.Sites[2].Status != "Up" ||
		result.Sites[2].HowLong != "2 days ago" || result.Sites[2].Name != "Test 3" {
		t.Error("Third site returned incorrect values")
	}

}
