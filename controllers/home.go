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

func (controller *homeController) get(rw http.ResponseWriter, req *http.Request) (int, error) {
	var sites database.Sites
	// Get active sites with no contacts.
	err := sites.GetSites(controller.DB, true, false)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Check if all of the active sites have a FirstPing that is not zero.  If so
	// then get the first ping from the database, if available.
	for i, site := range sites {
		if site.FirstPing.IsZero() {
			firstPing, err := site.GetFirstPing(controller.DB)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			err = site.UpdateSiteFirstPing(controller.DB, firstPing)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			sites[i].FirstPing = firstPing
		}
	}
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	messages := controller.authorizer.Messages(rw, req)
	vm := viewmodels.GetHomeViewModel(sites, isAuthenticated, user, messages)
	return http.StatusOK, controller.template.Execute(rw, vm)
}
