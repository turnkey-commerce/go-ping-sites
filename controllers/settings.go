package controllers

import (
	"database/sql"
	"net/http"
	"text/template"

	"github.com/apexskier/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type settingsController struct {
	DB         *sql.DB
	template   *template.Template
	authorizer httpauth.Authorizer
}

func (controller *settingsController) get(rw http.ResponseWriter, req *http.Request) {
	var sites database.Sites
	// Get all of the sites, including inactive ones, and the contacts.
	err := sites.GetSites(controller.DB, false, true)
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.GetSettingsViewModel(sites, isAuthenticated, user, err)
	controller.template.Execute(rw, vm)
}
