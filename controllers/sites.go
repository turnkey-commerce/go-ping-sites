package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"text/template"

	"github.com/apexskier/httpauth"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type sitesController struct {
	DB              *sql.DB
	detailsTemplate *template.Template
	editTemplate    *template.Template
	newTemplate     *template.Template
	authorizer      httpauth.Authorizer
}

func (controller *sitesController) getDetails(rw http.ResponseWriter, req *http.Request) (int, error) {
	vars := mux.Vars(req)
	siteID, err := strconv.ParseInt(vars["siteID"], 10, 64)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	site := new(database.Site)
	err = site.GetSite(controller.DB, siteID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.GetSiteDetailsViewModel(site, isAuthenticated, user)
	return http.StatusOK, controller.detailsTemplate.Execute(rw, vm)
}

//validateSiteForm checks the inputs for errors
func validateSiteForm(site *viewmodels.SitesEditViewModel) (valErrors map[string]string) {
	valErrors = make(map[string]string)

	_, err := govalidator.ValidateStruct(site)
	valErrors = govalidator.ErrorsByField(err)

	return valErrors
}

func mapSite(site *database.Site, formSite *viewmodels.SitesEditViewModel) {
	site.Name = formSite.Name
	site.IsActive = formSite.IsActive
	site.URL = formSite.URL
	site.PingIntervalSeconds = formSite.PingIntervalSeconds
	site.TimeoutSeconds = formSite.TimeoutSeconds
}
