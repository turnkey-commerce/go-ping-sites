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
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	messages := controller.authorizer.Messages(rw, req)
	vm := viewmodels.GetAboutViewModel(isAuthenticated, user, messages)
	controller.template.Execute(rw, vm)
}
