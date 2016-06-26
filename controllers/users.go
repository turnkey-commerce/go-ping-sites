package controllers

import (
	"net/http"
	"text/template"

	"github.com/apexskier/httpauth"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type usersController struct {
	getTemplate  *template.Template
	editTemplate *template.Template
	newTemplate  *template.Template
	authorizer   httpauth.Authorizer
	authBackend  httpauth.AuthBackend
	roles        map[string]httpauth.Role
}

func (controller *usersController) get(rw http.ResponseWriter, req *http.Request) (int, error) {
	// Get all of the users
	users, err := controller.authBackend.Users()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.GetUsersViewModel(users, isAuthenticated, user, err)
	return http.StatusOK, controller.getTemplate.Execute(rw, vm)
}

func (controller *usersController) editGet(rw http.ResponseWriter, req *http.Request) (int, error) {
	vars := mux.Vars(req)
	username := vars["username"]
	// Get the user to edit
	editUser, err := controller.authBackend.User(username)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	userEdit := new(viewmodels.UsersEditViewModel)
	userEdit.Email = editUser.Email
	userEdit.Role = editUser.Role
	userEdit.Username = editUser.Username
	vm := viewmodels.EditUserViewModel(userEdit, controller.roles, isAuthenticated, user, make(map[string]string))
	vm.CsrfField = csrf.TemplateField(req)
	return http.StatusOK, controller.editTemplate.Execute(rw, vm)
}

func (controller *usersController) editPost(rw http.ResponseWriter, req *http.Request) (int, error) {
	err := req.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	decoder := schema.NewDecoder()
	// Ignore unknown keys to prevent errors from the CSRF token.
	decoder.IgnoreUnknownKeys(true)
	formUser := new(viewmodels.UsersEditViewModel)
	err = decoder.Decode(formUser, req.PostForm)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	valErrors := validateUserForm(formUser, true)
	if len(valErrors) > 0 {
		isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
		vm := viewmodels.EditUserViewModel(formUser, controller.roles, isAuthenticated, user, valErrors)
		vm.CsrfField = csrf.TemplateField(req)
		return http.StatusOK, controller.editTemplate.Execute(rw, vm)
	}

	// Update the user.
	err = controller.authorizer.Update(rw, req, formUser.Username, formUser.Password, formUser.Email)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	http.Redirect(rw, req, "/settings/users", http.StatusSeeOther)
	return http.StatusSeeOther, nil
}

func (controller *usersController) newGet(rw http.ResponseWriter, req *http.Request) (int, error) {
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	userEdit := new(viewmodels.UsersEditViewModel)
	userEdit.Role = "user"
	vm := viewmodels.NewUserViewModel(userEdit, controller.roles, isAuthenticated, user, make(map[string]string))
	vm.CsrfField = csrf.TemplateField(req)
	return http.StatusOK, controller.newTemplate.Execute(rw, vm)
}

func (controller *usersController) newPost(rw http.ResponseWriter, req *http.Request) (int, error) {
	err := req.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	decoder := schema.NewDecoder()
	// Ignore unknown keys to prevent errors from the CSRF token.
	decoder.IgnoreUnknownKeys(true)
	formUser := new(viewmodels.UsersEditViewModel)
	err = decoder.Decode(formUser, req.PostForm)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	valErrors := validateUserForm(formUser, false)
	if len(valErrors) > 0 {
		isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
		vm := viewmodels.NewUserViewModel(formUser, controller.roles, isAuthenticated, user, valErrors)
		vm.CsrfField = csrf.TemplateField(req)
		return http.StatusOK, controller.newTemplate.Execute(rw, vm)
	}

	var user httpauth.UserData
	user.Username = formUser.Username
	user.Email = formUser.Email
	password := formUser.Password
	user.Role = formUser.Role
	err = controller.authorizer.Register(rw, req, user, password)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	http.Redirect(rw, req, "/settings/users", http.StatusSeeOther)
	return http.StatusSeeOther, nil
}

// validateUserForm checks the inputs for errors
func validateUserForm(user *viewmodels.UsersEditViewModel, allowMissingPassword bool) (valErrors map[string]string) {
	valErrors = make(map[string]string)

	_, err := govalidator.ValidateStruct(user)
	valErrors = govalidator.ErrorsByField(err)

	validatePassword(allowMissingPassword, user.Password, user.Password2, valErrors)

	return valErrors
}
