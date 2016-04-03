package controllers

import (
	"database/sql"
	"net/http"
	"text/template"

	"github.com/apexskier/httpauth"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type profileController struct {
	DB          *sql.DB
	template    *template.Template
	authorizer  httpauth.Authorizer
	authBackend httpauth.AuthBackend
}

func (controller *profileController) get(rw http.ResponseWriter, req *http.Request) (int, error) {
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	// Get the user to edit
	editUser, err := controller.authBackend.User(user.Username)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	userEdit := new(viewmodels.ProfileEditViewModel)
	userEdit.Email = editUser.Email
	vm := viewmodels.EditProfileViewModel(userEdit, isAuthenticated, user, make(map[string]string))
	return http.StatusOK, controller.template.Execute(rw, vm)
}

func (controller *profileController) post(rw http.ResponseWriter, req *http.Request) (int, error) {
	err := req.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	decoder := schema.NewDecoder()
	formUser := new(viewmodels.ProfileEditViewModel)
	err = decoder.Decode(formUser, req.PostForm)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	valErrors := validateProfileForm(formUser, true)
	if len(valErrors) > 0 {
		isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
		vm := viewmodels.EditProfileViewModel(formUser, isAuthenticated, user, valErrors)
		return http.StatusOK, controller.template.Execute(rw, vm)
	}

	// Update the user.
	err = controller.authorizer.Update(rw, req, "", formUser.Password, formUser.Email)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	http.Redirect(rw, req, "/", http.StatusSeeOther)
	return http.StatusSeeOther, nil
}

// validateUserForm checks the inputs for errors
func validateProfileForm(user *viewmodels.ProfileEditViewModel, allowMissingPassword bool) (valErrors map[string]string) {
	valErrors = make(map[string]string)

	_, err := govalidator.ValidateStruct(user)
	valErrors = govalidator.ErrorsByField(err)

	validatePassword(allowMissingPassword, user.Password, user.Password2, valErrors)

	return valErrors
}
