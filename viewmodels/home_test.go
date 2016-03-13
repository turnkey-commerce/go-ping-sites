package viewmodels_test

import (
	"testing"
	"time"

	"github.com/apexskier/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
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
	// Third site was up 2 days ago.
	twodaysAgo := now.Add(-48 * time.Hour)
	sites = append(sites, database.Site{Name: "Test 3", IsSiteUp: true, LastStatusChange: twodaysAgo})
	// Fourth site has no last status change but does have first ping three days ago.
	threedaysAgo := now.Add(-72 * time.Hour)
	sites = append(sites, database.Site{Name: "Test 4", IsSiteUp: true, FirstPing: threedaysAgo})

	result := viewmodels.GetHomeViewModel(sites, false, user, nil)
	if result.Nav.Active != "home" {
		t.Error("Home View Model Active returned incorrect value")
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
		result.Sites[2].HowLong != "2 days ago" || result.Sites[2].Name != "Test 3" ||
		result.Sites[2].HasNoStatusChanges {
		t.Error("Third site returned incorrect values")
	}

	if result.Sites[3].CSSClass != "success" || result.Sites[3].Status != "Up" ||
		result.Sites[3].HowLong != "3 days ago" || result.Sites[3].Name != "Test 4" ||
		!result.Sites[3].HasNoStatusChanges {
		t.Error("Fourth site returned incorrect values")
	}

	if result.HasSiteWithNoStatusChanges != true {
		t.Error("Should indicate has site with no status change.")
	}

}

// TestGetHomeViewModel tests the view model is created as expected.
func TestGetHomeViewModelWithNoStatusChanges(t *testing.T) {
	sites := database.Sites{}
	user := httpauth.UserData{}

	// First site has no last status change.
	sites = append(sites, database.Site{Name: "Test 1", IsSiteUp: true})
	// Second site was down 2 hours ago.
	now := time.Now()
	twoHoursAgo := now.Add(-2 * time.Hour)
	sites = append(sites, database.Site{Name: "Test 2", IsSiteUp: false, LastStatusChange: twoHoursAgo})
	// Third site was up 2 days ago.
	twodaysAgo := now.Add(-48 * time.Hour)
	sites = append(sites, database.Site{Name: "Test 3", IsSiteUp: true, LastStatusChange: twodaysAgo})

	result := viewmodels.GetHomeViewModel(sites, false, user, nil)

	if result.HasSiteWithNoStatusChanges != true {
		t.Error("Should indicate has site with no status change.")
	}
}

// TestGetHomeViewModel tests the view model is created as expected.
func TestGetHomeViewModelWithStatusChanges(t *testing.T) {
	sites := database.Sites{}
	user := httpauth.UserData{}

	// First site was up 2 hours ago.
	now := time.Now()
	twoHoursAgo := now.Add(-2 * time.Hour)
	sites = append(sites, database.Site{Name: "Test 1", IsSiteUp: true, LastStatusChange: twoHoursAgo})
	// Second site was down 2 hours ago.
	sites = append(sites, database.Site{Name: "Test 2", IsSiteUp: false, LastStatusChange: twoHoursAgo})
	// Third site was up 2 days ago.
	twodaysAgo := now.Add(-48 * time.Hour)
	sites = append(sites, database.Site{Name: "Test 3", IsSiteUp: true, LastStatusChange: twodaysAgo})

	result := viewmodels.GetHomeViewModel(sites, false, user, nil)

	if result.HasSiteWithNoStatusChanges == true {
		t.Error("Should NOT indicate has site with no status change.")
	}
}
