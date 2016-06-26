package viewmodels

import (
	"github.com/apexskier/httpauth"
	"html/template"
)

// UsersEditViewModel holds the required information about the site to choose for editing.
type UsersEditViewModel struct {
	Username  string `valid:"alphanum,required"`
	Email     string `valid:"email,required"`
	Role      string `valid:"-"`
	Password  string `valid:"-"`
	Password2 string `valid:"-"`
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
	Errors    map[string]string
	Title     string
	User      UsersEditViewModel
	Nav       NavViewModel
	Roles     map[string]httpauth.Role
	CsrfField template.HTML
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

// EditUserViewModel populates the items required by the user_edit.gohtml view
func EditUserViewModel(formUser *UsersEditViewModel, roles map[string]httpauth.Role,
	isAuthenticated bool, user httpauth.UserData, errors map[string]string) UserViewModel {
	nav := NavViewModel{
		Active:          "settings",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := UserViewModel{
		Title:  "Go Ping Sites - Settings - Edit User",
		Nav:    nav,
		Roles:  roles,
		Errors: errors,
	}

	userVM := new(UsersEditViewModel)
	userVM.Username = formUser.Username
	userVM.Email = formUser.Email
	userVM.Role = formUser.Role

	result.User = *userVM

	return result
}

// NewUserViewModel populates the items required by the user_new.gohtml view
func NewUserViewModel(formUser *UsersEditViewModel, roles map[string]httpauth.Role,
	isAuthenticated bool, user httpauth.UserData, errors map[string]string) UserViewModel {
	nav := NavViewModel{
		Active:          "settings",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := UserViewModel{
		Title:  "Go Ping Sites - Settings - New User",
		Nav:    nav,
		Roles:  roles,
		Errors: errors,
	}

	userVM := new(UsersEditViewModel)
	userVM.Username = formUser.Username
	userVM.Email = formUser.Email
	userVM.Role = formUser.Role
	userVM.Password = formUser.Password
	userVM.Password2 = formUser.Password2

	result.User = *userVM

	return result
}
