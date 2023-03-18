package controllers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAboutController(t *testing.T) {

	mockUserGetter := MockCurrentUserGetter{Username: "jules",
		FlashMessages: []string{"Log in to do that."}}

	req, _ := http.NewRequest("GET", "/about", nil)

	templates := populateFileTemplates("../templates")

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

// populateTemplates loads and parses all of the file templates in the templates directory
// This is for test only and is separate from the embedded templates in the main program.
func populateFileTemplates(templatePath string) *template.Template {
	result := template.New("templates")

	basePath := templatePath
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()

	templatePathsRaw, _ := templateFolder.Readdir(-1)
	templatePaths := new([]string)
	for _, pathInfo := range templatePathsRaw {
		if !pathInfo.IsDir() {
			*templatePaths = append(*templatePaths,
				basePath+"/"+pathInfo.Name())
		}
	}

	var funcMap = template.FuncMap{
		"displayBool":        displayBool,
		"displayActiveClass": displayActiveClass,
	}

	result.Funcs(funcMap).ParseFiles(*templatePaths...)
	return result
}
