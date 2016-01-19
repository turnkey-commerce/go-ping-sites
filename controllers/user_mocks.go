package controllers

import (
	"net/http"

	"github.com/apexskier/httpauth"
)

// MockCurrentUserGetter provides the base struct for the methods to mock the
// httpauth package. The properties on this struct can vary the results returned.
type MockCurrentUserGetter struct {
	FlashMessages []string
	Username      string
	UserError     error
}

// Messages is the mock of the Messasges method from the httpauth package.
func (m MockCurrentUserGetter) Messages(rw http.ResponseWriter, req *http.Request) []string {
	messages := m.FlashMessages
	return messages
}

// CurrentUser is the mock of the CurrentUser method from the httpauth package.
func (m MockCurrentUserGetter) CurrentUser(rw http.ResponseWriter, req *http.Request) (user httpauth.UserData, e error) {
	user = httpauth.UserData{Username: m.Username}
	return user, m.UserError
}
