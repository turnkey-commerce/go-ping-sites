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

type contactsController struct {
	DB           *sql.DB
	getTemplate  *template.Template
	editTemplate *template.Template
	newTemplate  *template.Template
	authorizer   httpauth.Authorizer
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

	vm := viewmodels.EditContactViewModel(contactEdit, isAuthenticated, user, make(map[string]string))
	return http.StatusOK, controller.editTemplate.Execute(rw, vm)
}

func (controller *contactsController) editPost(rw http.ResponseWriter, req *http.Request) (int, error) {
	err := req.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	decoder := schema.NewDecoder()
	formContact := new(viewmodels.ContactsEditViewModel)
	err = decoder.Decode(formContact, req.PostForm)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	valErrors := validateContactForm(formContact)
	if len(valErrors) > 0 {
		isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
		vm := viewmodels.NewContactViewModel(formContact, isAuthenticated, user, valErrors)
		return http.StatusOK, controller.newTemplate.Execute(rw, vm)
	}

	// Get the contact to edit
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
	http.Redirect(rw, req, "/settings/contacts", http.StatusSeeOther)
	return http.StatusSeeOther, nil
}

func (controller *contactsController) newGet(rw http.ResponseWriter, req *http.Request) (int, error) {
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	contactEdit := new(viewmodels.ContactsEditViewModel)
	contactEdit.EmailActive = false
	contactEdit.SmsActive = false
	vm := viewmodels.NewContactViewModel(contactEdit, isAuthenticated, user, make(map[string]string))
	return http.StatusOK, controller.newTemplate.Execute(rw, vm)
}

func (controller *contactsController) newPost(rw http.ResponseWriter, req *http.Request) (int, error) {
	err := req.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	decoder := schema.NewDecoder()
	formContact := new(viewmodels.ContactsEditViewModel)
	err = decoder.Decode(formContact, req.PostForm)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	valErrors := validateContactForm(formContact)
	if len(valErrors) > 0 {
		isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
		vm := viewmodels.NewContactViewModel(formContact, isAuthenticated, user, valErrors)
		return http.StatusOK, controller.newTemplate.Execute(rw, vm)
	}

	contact := database.Contact{}
	mapContacts(&contact, formContact)
	err = contact.CreateContact(controller.DB)
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

	return valErrors
}

func mapContacts(contact *database.Contact, formContact *viewmodels.ContactsEditViewModel) {
	contact.Name = formContact.Name
	contact.EmailAddress = formContact.EmailAddress
	contact.EmailActive = formContact.EmailActive
	contact.SmsNumber = formContact.SmsNumber
	contact.SmsActive = formContact.SmsActive
}
