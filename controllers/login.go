package controllers

import (
	"net/http"
	"text/template"

	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type loginController struct {
	template *template.Template
}

func (controller *loginController) get(w http.ResponseWriter, req *http.Request) {
	vm := viewmodels.GetLoginViewModel()
	controller.template.Execute(w, vm)
}
