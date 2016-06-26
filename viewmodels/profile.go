package viewmodels

import (
	"html/template"

	"github.com/apexskier/httpauth"
)

// ProfileEditViewModel holds the required information about the user to edit their profile.
type ProfileEditViewModel struct {
	Email     string `valid:"email,required"`
	Password  string `valid:"-"`
	Password2 string `valid:"-"`
	Username  string `valid:""`
}

// ProfileViewModel holds the view information for the profile.gohtml template
type ProfileViewModel struct {
	Errors    map[string]string
	Title     string
	User      ProfileEditViewModel
	Nav       NavViewModel
	CsrfField template.HTML
}

// EditProfileViewModel populates the items required by the profile.gohtml view
func EditProfileViewModel(formProfile *ProfileEditViewModel, isAuthenticated bool,
	user httpauth.UserData, errors map[string]string) ProfileViewModel {
	nav := NavViewModel{
		Active:          "profile",
		IsAuthenticated: isAuthenticated,
		User:            user,
	}

	result := ProfileViewModel{
		Title:  "Go Ping Sites - Profile",
		Nav:    nav,
		Errors: errors,
	}

	profileVM := new(ProfileEditViewModel)
	profileVM.Email = formProfile.Email
	profileVM.Password = formProfile.Password
	profileVM.Password2 = formProfile.Password2
	profileVM.Username = user.Username

	result.User = *profileVM

	return result
}
