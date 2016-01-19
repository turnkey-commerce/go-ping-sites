package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAboutController(t *testing.T) {

	mockUserGetter := MockCurrentUserGetter{Username: "jules",
		FlashMessages: []string{"Log in to do that."}}

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

	body := w.Body.String()

	if !strings.Contains(body, "Logged in as jules") {
		t.Error("Current user not rendered properly.")
	}

	if !strings.Contains(body, `<li class="active"><a href="/about">About</a></li>`) {
		t.Error("Navigation active tab not rendered as expected.")
	}

	if !strings.Contains(body, `<div class="alert alert-danger" role="alert">Log in to do that.</div>`) {
		t.Error("Flash message not rendered as expected.")
	}

	if !strings.Contains(body, `<title>Go Ping Sites - About</title>`) {
		t.Error("Flash message not rendered as expected.")
	}
}
