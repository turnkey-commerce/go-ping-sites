package controllers

import (
	"net/http"
	"text/template"

	"github.com/apexskier/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type usersController struct {
	template    *template.Template
	authorizer  httpauth.Authorizer
	authBackend httpauth.AuthBackend
}

func (controller *usersController) get(rw http.ResponseWriter, req *http.Request) {
	var users []httpauth.UserData
	// Get all of the users
	users, err := controller.authBackend.Users()
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.GetUsersViewModel(users, isAuthenticated, user, err)
	controller.template.Execute(rw, vm)
}
