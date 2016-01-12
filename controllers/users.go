package controllers

import (
	"net/http"
	"text/template"

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

func (controller *usersController) edit(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	// Get the user to edit
	editUser, err := controller.authBackend.User(username)
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.EditUserViewModel(editUser, controller.roles, isAuthenticated, user, err)
	controller.editTemplate.Execute(rw, vm)
}
