package controllers

import (
	"net/http"
	"text/template"

	"github.com/apexskier/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type aboutController struct {
	template   *template.Template
	authorizer httpauth.Authorizer
}

func (controller *aboutController) get(rw http.ResponseWriter, req *http.Request) {
	isAuthenticated := false
	user, authErr := controller.authorizer.CurrentUser(rw, req)
	if authErr == nil {
		isAuthenticated = true
	}
	vm := viewmodels.GetAboutViewModel(isAuthenticated, user)
	controller.template.Execute(rw, vm)
}
