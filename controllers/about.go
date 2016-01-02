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
	authErr := controller.authorizer.Authorize(rw, req, false)
	if authErr == nil {
		isAuthenticated = true
	}
	vm := viewmodels.GetAboutViewModel(isAuthenticated)
	controller.template.Execute(rw, vm)
}
