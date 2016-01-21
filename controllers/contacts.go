package controllers

import (
	"database/sql"
	"net/http"
	"text/template"

	"github.com/apexskier/httpauth"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type contactsController struct {
	DB          *sql.DB
	getTemplate *template.Template
	newTemplate *template.Template
	authorizer  httpauth.Authorizer
}

func (controller *contactsController) get(rw http.ResponseWriter, req *http.Request) {
	var contacts database.Contacts
	// Get contacts.
	err := contacts.GetContacts(controller.DB)
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.GetContactsViewModel(contacts, isAuthenticated, user, err)
	controller.getTemplate.Execute(rw, vm)
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
	// err = controller.authorizer.Register(rw, req, user, password)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Redirect(rw, req, "/settings/contacts/new", http.StatusSeeOther)
	// }
	http.Redirect(rw, req, "/settings/contacts", http.StatusSeeOther)
}

// validateUserForm checks the inputs for errors
func validateContactForm(contact *viewmodels.ContactsEditViewModel) (valErrors map[string]string) {
	valErrors = make(map[string]string)

	_, err := govalidator.ValidateStruct(contact)
	valErrors = govalidator.ErrorsByField(err)

	return valErrors
}
