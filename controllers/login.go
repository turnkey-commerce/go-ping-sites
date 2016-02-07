package controllers

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
	"github.com/apexskier/httpauth"
)

type loginController struct {
	template   *template.Template
	authorizer httpauth.Authorizer
}

// get creates the "/login" form.
func (controller *loginController) get(rw http.ResponseWriter, req *http.Request) {
	messages := controller.authorizer.Messages(rw, req)
	vm := viewmodels.GetLoginViewModel(messages)
	controller.template.Execute(rw, vm)
}

// post handles "/login" post requests.
func (controller *loginController) post(rw http.ResponseWriter, req *http.Request) {
	username := req.PostFormValue("username")
	password := req.PostFormValue("password")
	if err := controller.authorizer.Login(rw, req, username, password, "/"); err != nil && strings.Contains(err.Error(), "already authenticated") {
		http.Redirect(rw, req, "/", http.StatusSeeOther)
	} else if err != nil {
		http.Redirect(rw, req, "/login", http.StatusSeeOther)
	}
}
