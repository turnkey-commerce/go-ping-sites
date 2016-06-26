package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"text/template"

	"github.com/apexskier/httpauth"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/pinger"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
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
	// Get all of the contacts to display in the table.
	var contacts database.Contacts
	err = contacts.GetContacts(controller.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// Also get the site contacts to display in a table.
	err = site.GetSiteContacts(controller.DB, siteID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	selectedContacts := []int64{}
	for _, contact := range site.Contacts {
		selectedContacts = append(selectedContacts, contact.ContactID)
	}

	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	siteEdit := new(viewmodels.SitesEditViewModel)
	viewmodels.MapSiteDBtoVM(site, siteEdit)

	siteEdit.SelectedContacts = selectedContacts

	vm := viewmodels.EditSiteViewModel(siteEdit, contacts, isAuthenticated, user, make(map[string]string))
	vm.CsrfField = csrf.TemplateField(req)
	return http.StatusOK, controller.editTemplate.Execute(rw, vm)
}

func (controller *sitesController) editPost(rw http.ResponseWriter, req *http.Request) (int, error) {
	err := req.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	decoder := schema.NewDecoder()
	// Ignore unknown keys to prevent errors from the CSRF token.
	decoder.IgnoreUnknownKeys(true)
	formSite := new(viewmodels.SitesEditViewModel)
	err = decoder.Decode(formSite, req.PostForm)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	valErrors := validateSiteForm(formSite)
	if len(valErrors) > 0 {
		isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
		var contacts database.Contacts
		err = contacts.GetContacts(controller.DB)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		vm := viewmodels.EditSiteViewModel(formSite, contacts, isAuthenticated, user, valErrors)
		vm.CsrfField = csrf.TemplateField(req)
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

	//Loop selected ones first and if it's not already in the site then add it.
	for _, contactSelID := range formSite.SelectedContacts {
		if !int64InSlice(int64(contactSelID), formSite.SiteContacts) {
			err = site.AddContactToSite(controller.DB, contactSelID)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}

	// Loop existing site contacts and if it's not in the selected items then remove it.
	for _, contactSiteID := range formSite.SiteContacts {
		if !int64InSlice(int64(contactSiteID), formSite.SelectedContacts) {
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

	http.Redirect(rw, req, "/settings", http.StatusSeeOther)
	return http.StatusSeeOther, nil
}

func (controller *sitesController) newGet(rw http.ResponseWriter, req *http.Request) (int, error) {
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	siteNew := new(viewmodels.SitesEditViewModel)
	siteNew.IsActive = true
	// Get all of the contacts to display in the table.
	var contacts database.Contacts
	err := contacts.GetContacts(controller.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// These are strings in the ViewModel.
	siteNew.PingIntervalSeconds = "60"
	siteNew.TimeoutSeconds = "15"
	siteNew.SelectedContacts = []int64{}
	vm := viewmodels.NewSiteViewModel(siteNew, contacts, isAuthenticated, user, make(map[string]string))
	vm.CsrfField = csrf.TemplateField(req)
	return http.StatusOK, controller.newTemplate.Execute(rw, vm)
}

func (controller *sitesController) newPost(rw http.ResponseWriter, req *http.Request) (int, error) {
	err := req.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	decoder := schema.NewDecoder()
	// Ignore unknown keys to prevent errors from the CSRF token.
	decoder.IgnoreUnknownKeys(true)
	formSite := new(viewmodels.SitesEditViewModel)
	err = decoder.Decode(formSite, req.PostForm)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	valErrors := validateSiteForm(formSite)
	if len(valErrors) > 0 {
		isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
		var contacts database.Contacts
		err = contacts.GetContacts(controller.DB)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		vm := viewmodels.NewSiteViewModel(formSite, contacts, isAuthenticated, user, valErrors)
		vm.CsrfField = csrf.TemplateField(req)
		return http.StatusOK, controller.newTemplate.Execute(rw, vm)
	}

	site := database.Site{}
	viewmodels.MapSiteVMtoDB(formSite, &site)
	err = site.CreateSite(controller.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	//Add any selected contacts
	for _, contactSelID := range formSite.SelectedContacts {
		err = site.AddContactToSite(controller.DB, contactSelID)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	// Refresh the pinger with the changes.
	err = controller.pinger.UpdateSiteSettings()
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

func int64InSlice(a int64, list []int64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
