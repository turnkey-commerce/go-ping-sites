package controllers

import (
	"net/http"
	"text/template"

	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type aboutController struct {
	template *template.Template
}

func (controller *aboutController) get(w http.ResponseWriter, req *http.Request) {
	vm := viewmodels.GetAboutViewModel()
	controller.template.Execute(w, vm)
}
