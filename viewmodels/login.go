package viewmodels

// LoginViewModel holds the view information for the login.gohtml template
type LoginViewModel struct {
	Title    string
	Nav      NavViewModel
	Messages []string
}

// GetLoginViewModel populates the items required by the login.gohtml view
func GetLoginViewModel(messages []string) LoginViewModel {
	nav := NavViewModel{
		Active:          "login",
		IsAuthenticated: false,
	}

	result := LoginViewModel{
		Title:    "Go Ping Sites - Login",
		Nav:      nav,
		Messages: messages,
	}
	return result
}
