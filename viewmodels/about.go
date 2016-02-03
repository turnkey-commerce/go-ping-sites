package viewmodels

import "github.com/turnkey-commerce/httpauth"

// AboutViewModel holds the view information for the about.gohtml template
type AboutViewModel struct {
	Title    string
	Nav      NavViewModel
	Messages []string
}

// GetAboutViewModel populates the items required by the about.gohtml view
func GetAboutViewModel(isAuthenticated bool, user httpauth.UserData, messages []string) AboutViewModel {
	nav := NavViewModel{
		Active:          "about",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := AboutViewModel{
		Title:    "Go Ping Sites - About",
		Nav:      nav,
		Messages: messages,
	}
	return result
}
