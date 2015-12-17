package controllers

import (
	"net/http"
	"text/template"
)

type homeController struct {
	template *template.Template
}

func (controller *homeController) get(w http.ResponseWriter, req *http.Request) {
	controller.template.Execute(w, nil)
}
