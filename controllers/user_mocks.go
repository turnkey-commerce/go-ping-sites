package controllers

import (
	"net/http"

	"github.com/apexskier/httpauth"
)

type MockCurrentUserGetter struct {
}

func (m MockCurrentUserGetter) Messages(rw http.ResponseWriter, req *http.Request) []string {
	var messages []string
	return messages
}

func (m MockCurrentUserGetter) CurrentUser(rw http.ResponseWriter, req *http.Request) (user httpauth.UserData, e error) {
	user = httpauth.UserData{Username: "test"}
	return user, nil
}
