package controllers

import (
	"database/sql"
	"net/http"
	"text/template"

	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
	"github.com/turnkey-commerce/httpauth"
)

type settingsController struct {
	DB         *sql.DB
	template   *template.Template
	authorizer httpauth.Authorizer
}

func (controller *settingsController) get(rw http.ResponseWriter, req *http.Request) (int, error) {
	var sites database.Sites
	// Get all of the sites, including inactive ones, and the contacts.
	err := sites.GetSites(controller.DB, false, true)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.GetSettingsViewModel(sites, isAuthenticated, user, err)
	return http.StatusOK, controller.template.Execute(rw, vm)
}
