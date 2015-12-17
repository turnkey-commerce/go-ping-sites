package controllers

import (
	"net/http"
	"text/template"
)

type aboutController struct {
	template *template.Template
}

func (controller *aboutController) get(w http.ResponseWriter, req *http.Request) {
	controller.template.Execute(w, nil)
}
