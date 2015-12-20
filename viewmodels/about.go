package viewmodels

// AboutViewModel holds the view information for the about.html template
type AboutViewModel struct {
	Title  string
	Active string
}

// GetAboutViewModel populates the items required by the about.html view
func GetAboutViewModel() AboutViewModel {
	result := AboutViewModel{
		Title:  "Go Ping Sites - About",
		Active: "about",
	}
	return result
}
