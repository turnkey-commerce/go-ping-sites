package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"text/template"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"

	"github.com/apexskier/httpauth"
	"github.com/gorilla/mux"
	"github.com/turnkey-commerce/go-ping-sites/viewmodels"
)

type usersController struct {
	getTemplate  *template.Template
	editTemplate *template.Template
	newTemplate  *template.Template
	authorizer   httpauth.Authorizer
	authBackend  httpauth.AuthBackend
	roles        map[string]httpauth.Role
}

func (controller *usersController) get(rw http.ResponseWriter, req *http.Request) {
	// Get all of the users
	users, err := controller.authBackend.Users()
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.GetUsersViewModel(users, isAuthenticated, user, err)
	controller.getTemplate.Execute(rw, vm)
}

func (controller *usersController) editGet(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	// Get the user to edit
	editUser, err := controller.authBackend.User(username)
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.EditUserViewModel(editUser, controller.roles, isAuthenticated, user, err)
	controller.editTemplate.Execute(rw, vm)
}

func (controller *usersController) editPost(rw http.ResponseWriter, req *http.Request) {
	authErr := controller.authorizer.AuthorizeRole(rw, req, "admin", true)
	if authErr != nil {
		http.Redirect(rw, req, "/", http.StatusSeeOther)
	}
	password := req.PostFormValue("password")
	username := req.PostFormValue("username")
	role := req.PostFormValue("role")
	email := req.PostFormValue("email")

	// Get the user to edit
	var hash []byte
	editUser, err := controller.authBackend.User(username)
	if password != "" {
		hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		hash = editUser.Hash
	}

	newuser := httpauth.UserData{Username: username, Email: email, Hash: hash, Role: role}
	err = controller.authBackend.SaveUser(newuser)
	if err != nil {
		http.Redirect(rw, req, "/settings/users/"+username+"/edit", http.StatusSeeOther)
	}
	http.Redirect(rw, req, "/settings/users", http.StatusSeeOther)
}

func (controller *usersController) newGet(rw http.ResponseWriter, req *http.Request) {
	isAuthenticated, user := getCurrentUser(rw, req, controller.authorizer)
	vm := viewmodels.NewUserViewModel(controller.roles, isAuthenticated, user)
	controller.newTemplate.Execute(rw, vm)
}

func (controller *usersController) newPost(rw http.ResponseWriter, req *http.Request) {
	authErr := controller.authorizer.AuthorizeRole(rw, req, "admin", true)
	if authErr != nil {
		http.Redirect(rw, req, "/", http.StatusSeeOther)
	}

	var user httpauth.UserData
	user.Username = req.PostFormValue("username")
	user.Email = req.PostFormValue("email")
	password := req.PostFormValue("password")
	user.Role = req.PostFormValue("role")
	err := controller.authorizer.Register(rw, req, user, password)
	if err != nil {
		fmt.Println(err)
		http.Redirect(rw, req, "/settings/users/new", http.StatusSeeOther)
	}
	http.Redirect(rw, req, "/settings/users", http.StatusSeeOther)
}

// Message contains the inputs and any validation errors
type Message struct {
	Email    string
	Username string
	Password string
	Errors   map[string]string
}

// Validate checks the inputs for errors
func (msg *Message) Validate() bool {
	msg.Errors = make(map[string]string)

	if strings.TrimSpace(msg.Username) == "" {
		msg.Errors["Content"] = "Please provide a Username"
	}

	if utf8.RuneCountInString(strings.TrimSpace(msg.Password)) < 6 {
		msg.Errors["Content"] = "Password must be at least 6 characters in length"
	}

	re := regexp.MustCompile(".+@.+\\..+")
	matched := re.Match([]byte(msg.Email))
	if matched == false {
		msg.Errors["Email"] = "Please enter a valid email address"
	}

	return len(msg.Errors) == 0
}
