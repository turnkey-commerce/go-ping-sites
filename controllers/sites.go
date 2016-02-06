package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"text/template"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/pinger"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
	"github.com/turnkey-commerce/httpauth"
)

type sitesController struct {
	DB                     *sql.DB
	detailsTemplate        *template.Template
	editTemplate           *template.Template
	newTemplate            *template.Template
	changeContactsTemplate *template.Template
	authorizer             httpauth.Authorizer
	pinger                 *pinger.Pinger
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

	// Refresh the pinger with the changes.
	err = controller.pinger.UpdateSiteSettings()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(rw, req, "/settings", http.StatusSeeOther)
	return http.StatusSeeOther, nil
}

func (controller *sitesController) newGet(rw http.ResponseWriter, req *http.Request) (int, error) {
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	siteNew := new(viewmodels.SitesEditViewModel)
	siteNew.IsActive = true
	// These are strings in the ViewModel.
	siteNew.PingIntervalSeconds = "30"
	siteNew.TimeoutSeconds = "15"
	vm := viewmodels.NewSiteViewModel(siteNew, isAuthenticated, user, make(map[string]string))
	return http.StatusOK, controller.newTemplate.Execute(rw, vm)
}

func (controller *sitesController) newPost(rw http.ResponseWriter, req *http.Request) (int, error) {
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
		vm := viewmodels.NewSiteViewModel(formSite, isAuthenticated, user, valErrors)
		return http.StatusOK, controller.newTemplate.Execute(rw, vm)
	}

	site := database.Site{}
	viewmodels.MapSiteVMtoDB(formSite, &site)
	err = site.CreateSite(controller.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Refresh the pinger with the changes.
	err = controller.pinger.UpdateSiteSettings()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(rw, req, "/settings", http.StatusSeeOther)
	return http.StatusSeeOther, nil
}

func (controller *sitesController) editContactsGet(rw http.ResponseWriter, req *http.Request) (int, error) {
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
	// Get the contacts for the user to know the currently assigned ones.
	err = site.GetSiteContacts(controller.DB, siteID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// Get all of the contacts to display in the table.
	var contacts database.Contacts
	err = contacts.GetContacts(controller.DB)

	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.SiteChangeContactsViewModel(site, contacts, isAuthenticated, user)
	return http.StatusOK, controller.changeContactsTemplate.Execute(rw, vm)
}

func (controller *sitesController) editContactsPost(rw http.ResponseWriter, req *http.Request) (int, error) {
	err := req.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	decoder := schema.NewDecoder()
	formContacts := new(viewmodels.SiteContactsSelectedViewModel)
	err = decoder.Decode(formContacts, req.PostForm)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	site := new(database.Site)
	err = site.GetSite(controller.DB, formContacts.SiteID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	//Loop selected ones first and if it's not already in the site then add it.
	for _, contactSelID := range formContacts.SelectedContacts {
		if !int64InSlice(int64(contactSelID), formContacts.SiteContacts) {
			err = site.AddContactToSite(controller.DB, contactSelID)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}

	// Loop existing site contacts and if it's not in the selected items then remove it.
	for _, contactSiteID := range formContacts.SiteContacts {
		if !int64InSlice(int64(contactSiteID), formContacts.SelectedContacts) {
			err = site.RemoveContactFromSite(controller.DB, contactSiteID)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}

	// Refresh the pinger with the changes.
	err = controller.pinger.UpdateSiteSettings()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(rw, req, "/settings/sites/"+strconv.FormatInt(site.SiteID, 10), http.StatusSeeOther)
	return http.StatusOK, nil
}

//validateSiteForm checks the inputs for errors
func validateSiteForm(site *viewmodels.SitesEditViewModel) (valErrors map[string]string) {
	valErrors = make(map[string]string)
	_, err := govalidator.ValidateStruct(site)
	valErrors = govalidator.ErrorsByField(err)
	return valErrors
}

func int64InSlice(a int64, list []int64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
