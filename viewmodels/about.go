package viewmodels

// AboutViewModel holds the view information for the about.gohtml template
type AboutViewModel struct {
	Title string
	Nav   NavViewModel
}

// GetAboutViewModel populates the items required by the about.gohtml view
func GetAboutViewModel(isAuthenticated bool) AboutViewModel {
	nav := NavViewModel{
		Active:          "about",
		IsAuthenticated: isAuthenticated,
	}

	result := AboutViewModel{
		Title: "Go Ping Sites - About",
		Nav:   nav,
	}
	return result
}
