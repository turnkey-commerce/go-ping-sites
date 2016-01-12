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

// UserViewModel holds the view information for the user_edit.gohtml template
type UserViewModel struct {
	Error error
	Title string
	User  UsersEditViewModel
	Nav   NavViewModel
	Roles map[string]httpauth.Role
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

// EditUserViewModel populates the items required by the settings.gohtml view
func EditUserViewModel(editUser httpauth.UserData, roles map[string]httpauth.Role, isAuthenticated bool, user httpauth.UserData, err error) UserViewModel {
	nav := NavViewModel{
		Active:          "settings",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := UserViewModel{
		Title: "Go Ping Sites - Settings - Edit User",
		Nav:   nav,
		Roles: roles,
	}

	userVM := new(UsersEditViewModel)
	userVM.Username = editUser.Username
	userVM.Email = editUser.Email
	userVM.Role = editUser.Role

	result.User = *userVM

	return result
}
