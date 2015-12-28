package viewmodels

// LoginViewModel holds the view information for the login.gohtml template
type LoginViewModel struct {
	Title  string
	Active string
}

// GetLoginViewModel populates the items required by the login.gohtml view
func GetLoginViewModel() LoginViewModel {
	result := LoginViewModel{
		Title:  "Go Ping Sites - Login",
		Active: "login",
	}
	return result
}
