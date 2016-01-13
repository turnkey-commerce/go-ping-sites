package controllers

import (
	"fmt"
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"

	"github.com/apexskier/httpauth"
	"github.com/gorilla/mux"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type usersController struct {
	getTemplate  *template.Template
	editTemplate *template.Template
	authorizer   httpauth.Authorizer
	authBackend  httpauth.AuthBackend
	roles        map[string]httpauth.Role
}

func (controller *usersController) get(rw http.ResponseWriter, req *http.Request) {
	// Get all of the users
	users, err := controller.authBackend.Users()
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.GetUsersViewModel(users, isAuthenticated, user, err)
	controller.getTemplate.Execute(rw, vm)
}

func (controller *usersController) editGet(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	// Get the user to edit
	editUser, err := controller.authBackend.User(username)
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.EditUserViewModel(editUser, controller.roles, isAuthenticated, user, err)
	controller.editTemplate.Execute(rw, vm)
}

func (controller *usersController) editPost(rw http.ResponseWriter, req *http.Request) {
	authErr := controller.authorizer.AuthorizeRole(rw, req, "admin", true)
	if authErr != nil {
		http.Redirect(rw, req, "/", http.StatusSeeOther)
	}
	password := req.PostFormValue("password")
	username := req.PostFormValue("username")
	role := req.PostFormValue("role")
	email := req.PostFormValue("email")

	// Get the user to edit
	var hash []byte
	editUser, err := controller.authBackend.User(username)
	if password != "" {
		hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		hash = editUser.Hash
	}

	newuser := httpauth.UserData{Username: username, Email: email, Hash: hash, Role: role}
	err = controller.authBackend.SaveUser(newuser)
	if err != nil {
		http.Redirect(rw, req, "/settings/users/"+username+"/edit", http.StatusSeeOther)
	}
	http.Redirect(rw, req, "/settings/users", http.StatusSeeOther)
}
