package viewmodels

// AboutViewModel holds the view information for the about.gohtml template
type AboutViewModel struct {
	Title  string
	Active string
}

// GetAboutViewModel populates the items required by the about.gohtml view
func GetAboutViewModel() AboutViewModel {
	result := AboutViewModel{
		Title:  "Go Ping Sites - About",
		Active: "about",
	}
	return result
}
