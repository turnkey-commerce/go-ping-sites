package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"text/template"

	"github.com/apexskier/httpauth"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
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
	// Also get the contacts to display in a table.
	err = site.GetSiteContacts(controller.DB, siteID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.GetSiteDetailsViewModel(site, isAuthenticated, user)
	return http.StatusOK, controller.detailsTemplate.Execute(rw, vm)
}

func (controller *sitesController) editGet(rw http.ResponseWriter, req *http.Request) (int, error) {
	vars := mux.Vars(req)
	siteID, err := strconv.ParseInt(vars["siteID"], 10, 64)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// Get the site to edit
	site := new(database.Site)
	err = site.GetSite(controller.DB, siteID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	siteEdit := new(viewmodels.SitesEditViewModel)
	viewmodels.MapSiteDBtoVM(site, siteEdit)

	vm := viewmodels.EditSiteViewModel(siteEdit, isAuthenticated, user, make(map[string]string))
	return http.StatusOK, controller.editTemplate.Execute(rw, vm)
}

func (controller *sitesController) editPost(rw http.ResponseWriter, req *http.Request) (int, error) {
	err := req.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	decoder := schema.NewDecoder()
	formSite := new(viewmodels.SitesEditViewModel)
	err = decoder.Decode(formSite, req.PostForm)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	valErrors := validateSiteForm(formSite)
	if len(valErrors) > 0 {
		isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
		vm := viewmodels.EditSiteViewModel(formSite, isAuthenticated, user, valErrors)
		return http.StatusOK, controller.editTemplate.Execute(rw, vm)
	}

	// Get the site to edit
	site := new(database.Site)
	err = site.GetSite(controller.DB, formSite.SiteID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = viewmodels.MapSiteVMtoDB(formSite, site)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = site.UpdateSite(controller.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	http.Redirect(rw, req, "/settings", http.StatusSeeOther)
	return http.StatusSeeOther, nil
}

//validateSiteForm checks the inputs for errors
func validateSiteForm(site *viewmodels.SitesEditViewModel) (valErrors map[string]string) {
	valErrors = make(map[string]string)
	_, err := govalidator.ValidateStruct(site)
	valErrors = govalidator.ErrorsByField(err)
	return valErrors
}
