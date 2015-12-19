package controllers

import (
	"database/sql"
	"net/http"
	"text/template"

	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type homeController struct {
	DB       *sql.DB
	template *template.Template
}

func (controller *homeController) get(w http.ResponseWriter, req *http.Request) {
	vm := viewmodels.GetHomeViewModel(controller.DB)
	controller.template.Execute(w, vm)
}
