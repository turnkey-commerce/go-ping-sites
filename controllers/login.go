package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/apexskier/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type loginController struct {
	template   *template.Template
	authorizer httpauth.Authorizer
}

// get creates the "/login" form.
func (controller *loginController) get(rw http.ResponseWriter, req *http.Request) {
	vm := viewmodels.GetLoginViewModel()
	controller.template.Execute(rw, vm)
}

// post handles "/login" post requests.
func (controller *loginController) post(rw http.ResponseWriter, req *http.Request) {
	username := req.PostFormValue("username")
	password := req.PostFormValue("password")
	controller.authorizer.Logout(rw, req)
	if err := controller.authorizer.Login(rw, req, username, password, "/"); err != nil && strings.Contains(err.Error(), "already authenticated") {
		http.Redirect(rw, req, "/", http.StatusSeeOther)
	} else if err != nil {
		fmt.Println(err.Error())
		http.Redirect(rw, req, "/login", http.StatusSeeOther)
	}
}
