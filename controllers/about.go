package controllers

import (
	"html/template"
	"net/http"

	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type aboutController struct {
	template   *template.Template
	authorizer CurrentUserGetter
	version    string
}

func (controller *aboutController) get(rw http.ResponseWriter, req *http.Request) {
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	messages := controller.authorizer.Messages(rw, req)
	vm := viewmodels.GetAboutViewModel(isAuthenticated, user, messages, controller.version)
	controller.template.Execute(rw, vm)
}
