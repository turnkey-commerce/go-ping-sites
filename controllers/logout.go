package controllers

import (
	"net/http"

	"github.com/apexskier/httpauth"
)

type logoutController struct {
	authorizer httpauth.Authorizer
}

// get executes the logout.
func (controller *logoutController) get(rw http.ResponseWriter, req *http.Request) {
	controller.authorizer.Logout(rw, req)
	http.Redirect(rw, req, "/", http.StatusSeeOther)
}
