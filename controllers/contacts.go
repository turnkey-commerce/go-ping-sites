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

func (controller *contactsController) get(rw http.ResponseWriter, req *http.Request) {
	var contacts database.Contacts
	// Get contacts.
	err := contacts.GetContacts(controller.DB)
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.GetContactsViewModel(contacts, isAuthenticated, user, err)
	controller.getTemplate.Execute(rw, vm)
}

func (controller *contactsController) editGet(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	contactID, err := strconv.ParseInt(vars["contactID"], 10, 64)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get the contact to edit
	contact := new(database.Contact)
	err = contact.GetContact(controller.DB, contactID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
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
	controller.editTemplate.Execute(rw, vm)
}

func (controller *contactsController) newGet(rw http.ResponseWriter, req *http.Request) {
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	contactEdit := new(viewmodels.ContactsEditViewModel)
	contactEdit.EmailActive = false
	contactEdit.SmsActive = false
	vm := viewmodels.NewContactViewModel(contactEdit, isAuthenticated, user, make(map[string]string))
	controller.newTemplate.Execute(rw, vm)
}

func (controller *contactsController) newPost(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	decoder := schema.NewDecoder()
	formContact := new(viewmodels.ContactsEditViewModel)
	err = decoder.Decode(formContact, req.PostForm)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	valErrors := validateContactForm(formContact)
	if len(valErrors) > 0 {
		isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
		vm := viewmodels.NewContactViewModel(formContact, isAuthenticated, user, valErrors)
		controller.newTemplate.Execute(rw, vm)
		return
	}

	var contact database.Contact
	contact.Name = formContact.Name
	contact.EmailAddress = formContact.EmailAddress
	contact.EmailActive = formContact.EmailActive
	contact.SmsNumber = formContact.SmsNumber
	contact.SmsActive = formContact.SmsActive
	err = contact.CreateContact(controller.DB)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, req, "/settings/contacts", http.StatusSeeOther)
}

// validateUserForm checks the inputs for errors
func validateContactForm(contact *viewmodels.ContactsEditViewModel) (valErrors map[string]string) {
	valErrors = make(map[string]string)

	_, err := govalidator.ValidateStruct(contact)
	valErrors = govalidator.ErrorsByField(err)

	return valErrors
}
