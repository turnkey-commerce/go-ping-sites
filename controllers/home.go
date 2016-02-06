package controllers

import (
	"database/sql"
	"net/http"
	"text/template"

	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
	"github.com/turnkey-commerce/httpauth"
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
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	messages := controller.authorizer.Messages(rw, req)
	vm := viewmodels.GetHomeViewModel(sites, isAuthenticated, user, messages, err)
	controller.template.Execute(rw, vm)
}
