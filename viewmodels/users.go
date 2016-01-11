package viewmodels

import "github.com/apexskier/httpauth"

// UsersEditViewModel holds the required information about the site to choose for editing.
type UsersEditViewModel struct {
	Username string
	Email    string
	Role     string
}

// UsersViewModel holds the view information for the users.gohtml template
type UsersViewModel struct {
	Error error
	Title string
	Users []UsersEditViewModel
	Nav   NavViewModel
}

// GetUsersViewModel populates the items required by the settings.gohtml view
func GetUsersViewModel(users []httpauth.UserData, isAuthenticated bool, user httpauth.UserData, err error) UsersViewModel {
	nav := NavViewModel{
		Active:          "settings",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := UsersViewModel{
		Title: "Go Ping Sites - Settings - Users",
		Nav:   nav,
	}

	for _, user := range users {
		userVM := new(UsersEditViewModel)
		userVM.Username = user.Username
		userVM.Email = user.Email
		userVM.Role = user.Role
		result.Users = append(result.Users, *userVM)
	}

	return result
}
