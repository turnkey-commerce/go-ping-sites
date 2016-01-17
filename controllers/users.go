package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"

	"github.com/apexskier/httpauth"
	"github.com/asaskevich/govalidator"
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

func (controller *usersController) get(rw http.ResponseWriter, req *http.Request) {
	// Get all of the users
	users, err := controller.authBackend.Users()
	// TODO: Replace with better error handling.
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.GetUsersViewModel(users, isAuthenticated, user, err)
	controller.getTemplate.Execute(rw, vm)
}

func (controller *usersController) editGet(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	// Get the user to edit
	editUser, err := controller.authBackend.User(username)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	userEdit := new(viewmodels.UsersEditViewModel)
	userEdit.Email = editUser.Email
	userEdit.Role = editUser.Role
	userEdit.Username = editUser.Username
	vm := viewmodels.EditUserViewModel(userEdit, controller.roles, isAuthenticated, user, make(map[string]string))
	controller.editTemplate.Execute(rw, vm)
}

func (controller *usersController) editPost(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	decoder := schema.NewDecoder()
	formUser := new(viewmodels.UsersEditViewModel)
	err = decoder.Decode(formUser, req.PostForm)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	valErrors := validateUserForm(formUser, true)
	if len(valErrors) > 0 {
		isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
		vm := viewmodels.EditUserViewModel(formUser, controller.roles, isAuthenticated, user, valErrors)
		controller.editTemplate.Execute(rw, vm)
		return
	}

	// Get the user to edit
	var hash []byte
	editUser, err := controller.authBackend.User(formUser.Username)
	if formUser.Password != "" {
		hash, err = bcrypt.GenerateFromPassword([]byte(formUser.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		hash = editUser.Hash
	}

	newuser := httpauth.UserData{Username: formUser.Username, Email: formUser.Email, Hash: hash, Role: formUser.Role}
	err = controller.authBackend.SaveUser(newuser)
	if err != nil {
		http.Redirect(rw, req, "/settings/users/"+formUser.Username+"/edit", http.StatusSeeOther)
	}
	http.Redirect(rw, req, "/settings/users", http.StatusSeeOther)
}

func (controller *usersController) newGet(rw http.ResponseWriter, req *http.Request) {
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	userEdit := new(viewmodels.UsersEditViewModel)
	userEdit.Role = "user"
	vm := viewmodels.NewUserViewModel(userEdit, controller.roles, isAuthenticated, user, make(map[string]string))
	controller.newTemplate.Execute(rw, vm)
}

func (controller *usersController) newPost(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	decoder := schema.NewDecoder()
	formUser := new(viewmodels.UsersEditViewModel)
	err = decoder.Decode(formUser, req.PostForm)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	valErrors := validateUserForm(formUser, false)
	if len(valErrors) > 0 {
		isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
		vm := viewmodels.NewUserViewModel(formUser, controller.roles, isAuthenticated, user, valErrors)
		controller.newTemplate.Execute(rw, vm)
		return
	}

	var user httpauth.UserData
	user.Username = formUser.Username
	user.Email = formUser.Email
	password := formUser.Password
	user.Role = formUser.Role
	err = controller.authorizer.Register(rw, req, user, password)
	if err != nil {
		fmt.Println(err)
		http.Redirect(rw, req, "/settings/users/new", http.StatusSeeOther)
	}
	http.Redirect(rw, req, "/settings/users", http.StatusSeeOther)
}

// validateUserForm checks the inputs for errors
func validateUserForm(user *viewmodels.UsersEditViewModel, allowMissingPassword bool) (valErrors map[string]string) {
	valErrors = make(map[string]string)

	_, err := govalidator.ValidateStruct(user)
	valErrors = govalidator.ErrorsByField(err)

	if !allowMissingPassword && utf8.RuneCountInString(strings.TrimSpace(user.Password)) < 6 {
		valErrors["Password"] = "Password must be at least 6 characters in length"
	} else if allowMissingPassword && utf8.RuneCountInString(strings.TrimSpace(user.Password)) < 6 &&
		utf8.RuneCountInString(strings.TrimSpace(user.Password)) > 0 {
		valErrors["Password"] = "Password must be at least 6 characters in length"
	}

	return valErrors
}
