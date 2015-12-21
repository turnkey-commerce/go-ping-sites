package controllers

import (
	"database/sql"
	"net/http"
	"text/template"

	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type homeController struct {
	DB       *sql.DB
	template *template.Template
}

func (controller *homeController) get(w http.ResponseWriter, req *http.Request) {
	var sites database.Sites
	err := sites.GetActiveSites(controller.DB)
	vm := viewmodels.GetHomeViewModel(sites, err)
	controller.template.Execute(w, vm)
}
