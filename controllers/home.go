package controllers

import (
	"database/sql"
	"net/http"
	"text/template"

	"github.com/apexskier/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type homeController struct {
	DB         *sql.DB
	template   *template.Template
	authorizer httpauth.Authorizer
}

func (controller *homeController) get(rw http.ResponseWriter, req *http.Request) {
	var sites database.Sites
	// Get active sites with no contacts.
	err := sites.GetSites(controller.DB, true, false)
	isAuthenticated := false
	user, authErr := controller.authorizer.CurrentUser(rw, req)
	if authErr == nil {
		isAuthenticated = true
	}
	messages := controller.authorizer.Messages(rw, req)
	vm := viewmodels.GetHomeViewModel(sites, isAuthenticated, user, messages, err)
	controller.template.Execute(rw, vm)
}
