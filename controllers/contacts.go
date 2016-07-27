package controllers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"

	"github.com/apexskier/httpauth"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/pinger"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type contactsController struct {
	DB             *sql.DB
	getTemplate    *template.Template
	editTemplate   *template.Template
	newTemplate    *template.Template
	deleteTemplate *template.Template
	authorizer     httpauth.Authorizer
	pinger         *pinger.Pinger
}

func (controller *contactsController) get(rw http.ResponseWriter, req *http.Request) (int, error) {
	var contacts database.Contacts
	// Get contacts.
	err := contacts.GetContacts(controller.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.GetContactsViewModel(contacts, isAuthenticated, user, err)
	return http.StatusOK, controller.getTemplate.Execute(rw, vm)
}

func (controller *contactsController) editGet(rw http.ResponseWriter, req *http.Request) (int, error) {
	vars := mux.Vars(req)
	contactID, err := strconv.ParseInt(vars["contactID"], 10, 64)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// Get the contact to edit
	contact := new(database.Contact)
	err = contact.GetContact(controller.DB, contactID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	contactEdit := new(viewmodels.ContactsEditViewModel)
	contactEdit.Name = contact.Name
	contactEdit.ContactID = contact.ContactID
	contactEdit.EmailAddress = contact.EmailAddress
	contactEdit.SmsNumber = contact.SmsNumber
	contactEdit.EmailActive = contact.EmailActive
	contactEdit.SmsActive = contact.SmsActive
	contactEdit.SelectedSites, err = getContactSiteIDs(controller, contact)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	sites, errGet := getAllSites(controller)
	if errGet != nil {
		return http.StatusInternalServerError, errGet
	}

	vm := viewmodels.EditContactViewModel(contactEdit, sites, isAuthenticated, user, make(map[string]string))
	vm.CsrfField = csrf.TemplateField(req)
	return http.StatusOK, controller.editTemplate.Execute(rw, vm)
}

func (controller *contactsController) editPost(rw http.ResponseWriter, req *http.Request) (int, error) {
	err := req.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	decoder := schema.NewDecoder()
	// Ignore unknown keys to prevent errors from the CSRF token.
	decoder.IgnoreUnknownKeys(true)
	formContact := new(viewmodels.ContactsEditViewModel)
	err = decoder.Decode(formContact, req.PostForm)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	valErrors := validateContactForm(formContact)
	if len(valErrors) > 0 {
		isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
		sites, errGet := getAllSites(controller)
		if errGet != nil {
			return http.StatusInternalServerError, err
		}
		vm := viewmodels.EditContactViewModel(formContact, sites, isAuthenticated,
			user, valErrors)
		vm.CsrfField = csrf.TemplateField(req)
		return http.StatusOK, controller.editTemplate.Execute(rw, vm)
	}

	// Get the contact to update
	contact := new(database.Contact)
	err = contact.GetContact(controller.DB, formContact.ContactID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	mapContacts(contact, formContact)
	err = contact.UpdateContact(controller.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	contactSiteIDS, getErr := getContactSiteIDs(controller, contact)
	if getErr != nil {
		return http.StatusInternalServerError, getErr
	}
	//Loop selected ones first and if it's not already in the site then add it.
	for _, siteSelID := range formContact.SelectedSites {
		if !int64InSlice(int64(siteSelID), contactSiteIDS) {
			err = addContactToSite(controller, contact.ContactID, siteSelID)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}

	// Loop existing contact sites and if it's not in the selected items then remove it.
	for _, contactSiteID := range contactSiteIDS {
		if !int64InSlice(int64(contactSiteID), formContact.SelectedSites) {
			err = removeContactFromSite(controller, contact.ContactID, contactSiteID)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}

	// Refresh the pinger with the changes.
	// TODO: Check whether this contact is associated with any active site first.
	err = controller.pinger.UpdateSiteSettings()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(rw, req, "/settings/contacts", http.StatusSeeOther)
	return http.StatusSeeOther, nil
}

func (controller *contactsController) newGet(rw http.ResponseWriter, req *http.Request) (int, error) {
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	contactEdit := new(viewmodels.ContactsEditViewModel)
	contactEdit.EmailActive = false
	contactEdit.SmsActive = false
	sites, err := getAllSites(controller)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	vm := viewmodels.NewContactViewModel(contactEdit, sites, true,
		isAuthenticated, user, make(map[string]string))
	vm.CsrfField = csrf.TemplateField(req)
	return http.StatusOK, controller.newTemplate.Execute(rw, vm)
}

func (controller *contactsController) newPost(rw http.ResponseWriter, req *http.Request) (int, error) {
	err := req.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	decoder := schema.NewDecoder()
	// Ignore unknown keys to prevent errors from the CSRF token.
	decoder.IgnoreUnknownKeys(true)
	formContact := new(viewmodels.ContactsEditViewModel)
	err = decoder.Decode(formContact, req.PostForm)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	valErrors := validateContactForm(formContact)
	if len(valErrors) > 0 {
		isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
		sites, errGet := getAllSites(controller)
		if errGet != nil {
			return http.StatusInternalServerError, err
		}
		vm := viewmodels.NewContactViewModel(formContact, sites, false,
			isAuthenticated, user, valErrors)
		vm.CsrfField = csrf.TemplateField(req)
		return http.StatusOK, controller.newTemplate.Execute(rw, vm)
	}

	contact := database.Contact{}
	mapContacts(&contact, formContact)
	err = contact.CreateContact(controller.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	//Add contact to any selected sites
	for _, siteSelID := range formContact.SelectedSites {
		err = addContactToSite(controller, contact.ContactID, siteSelID)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	// Refresh the pinger with the changes.
	// TODO: Check whether this contact has been added to any site first.
	err = controller.pinger.UpdateSiteSettings()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(rw, req, "/settings/contacts", http.StatusSeeOther)
	return http.StatusSeeOther, nil
}

func (controller *contactsController) deleteGet(rw http.ResponseWriter, req *http.Request) (int, error) {
	vars := mux.Vars(req)
	contactID, err := strconv.ParseInt(vars["contactID"], 10, 64)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// Get the contact to edit
	contact := new(database.Contact)
	err = contact.GetContact(controller.DB, contactID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	contactDelete := new(viewmodels.ContactsEditViewModel)
	contactDelete.Name = contact.Name
	contactDelete.ContactID = contact.ContactID
	contactDelete.EmailAddress = contact.EmailAddress
	var noSites = []database.Site{}
	vm := viewmodels.EditContactViewModel(contactDelete, noSites, isAuthenticated,
		user, make(map[string]string))
	vm.CsrfField = csrf.TemplateField(req)
	return http.StatusOK, controller.deleteTemplate.Execute(rw, vm)
}

func (controller *contactsController) deletePost(rw http.ResponseWriter, req *http.Request) (int, error) {
	err := req.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	decoder := schema.NewDecoder()
	// Ignore unknown keys to prevent errors from the CSRF token.
	decoder.IgnoreUnknownKeys(true)
	formContact := new(viewmodels.ContactsEditViewModel)
	err = decoder.Decode(formContact, req.PostForm)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	valErrors := validateContactForm(formContact)
	if len(valErrors) > 0 {
		isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
		var noSites = []database.Site{}
		vm := viewmodels.EditContactViewModel(formContact, noSites, isAuthenticated,
			user, valErrors)
		vm.CsrfField = csrf.TemplateField(req)
		return http.StatusOK, controller.deleteTemplate.Execute(rw, vm)
	}

	// Get the contact to delete
	contact := new(database.Contact)
	err = contact.GetContact(controller.DB, formContact.ContactID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	mapContacts(contact, formContact)
	err = contact.DeleteContact(controller.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Refresh the pinger with the changes.
	// TODO: Check whether this contact is associated with any active site first.
	err = controller.pinger.UpdateSiteSettings()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.Redirect(rw, req, "/settings/contacts", http.StatusSeeOther)
	return http.StatusSeeOther, nil
}

// validateUserForm checks the inputs for errors
func validateContactForm(contact *viewmodels.ContactsEditViewModel) (valErrors map[string]string) {
	valErrors = make(map[string]string)

	_, err := govalidator.ValidateStruct(contact)
	valErrors = govalidator.ErrorsByField(err)

	validateContact(contact, valErrors)

	return valErrors
}

func mapContacts(contact *database.Contact, formContact *viewmodels.ContactsEditViewModel) {
	contact.Name = formContact.Name
	contact.EmailAddress = formContact.EmailAddress
	contact.EmailActive = formContact.EmailActive
	contact.SmsNumber = formContact.SmsNumber
	contact.SmsActive = formContact.SmsActive
}

func getAllSites(controller *contactsController) (database.Sites, error) {
	// Get all of the sites to display in the sites-to-assign table.
	var sites database.Sites
	err := sites.GetSites(controller.DB, false, false)
	if err != nil {
		return nil, err
	}
	return sites, nil
}

func getContactSiteIDs(controller *contactsController, contact *database.Contact) ([]int64, error) {
	// Get the site ID's for a given contact
	siteIDs := []int64{}
	err := contact.GetContactSites(controller.DB)
	if err != nil {
		return nil, err
	}
	for _, site := range contact.Sites {
		siteIDs = append(siteIDs, site.SiteID)
	}
	return siteIDs, nil
}

func addContactToSite(controller *contactsController, contactID int64, siteID int64) error {
	var site database.Site
	err := site.GetSite(controller.DB, siteID)
	if err != nil {
		return err
	}
	err = site.AddContactToSite(controller.DB, contactID)
	if err != nil {
		return err
	}
	return nil
}

func removeContactFromSite(controller *contactsController, contactID int64, siteID int64) error {
	var site database.Site
	err := site.GetSite(controller.DB, siteID)
	if err != nil {
		return err
	}
	err = site.RemoveContactFromSite(controller.DB, contactID)
	if err != nil {
		return err
	}
	return nil
}
