package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAboutController(t *testing.T) {

	mockUserGetter := MockCurrentUserGetter{}
	req, err := http.NewRequest("GET", "http://example.com/about", nil)
	if err != nil {
		t.Fatal(err)
	}

	templates := PopulateTemplates("../templates")

	ac := new(aboutController)
	ac.template = templates.Lookup("about.gohtml")
	ac.authorizer = mockUserGetter

	w := httptest.NewRecorder()
	ac.get(w, req)

	if !strings.Contains(w.Body.String(), "Logged in as test") {
		t.Error("Current user not rendered properly.")
	}

	if !strings.Contains(w.Body.String(), `<li class="active"><a href="/about">About</a></li>`) {
		t.Error("Navigation active tab not rendered as expected.")
	}
}
