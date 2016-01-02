package controllers

import (
	"bufio"
	"database/sql"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/apexskier/httpauth"
	"github.com/gorilla/mux"
)

// Register the handlers for a given route.
func Register(db *sql.DB, authorizer httpauth.Authorizer, templates *template.Template) {
	router := mux.NewRouter()

	hc := new(homeController)
	hc.template = templates.Lookup("home.gohtml")
	hc.authorizer = authorizer
	hc.DB = db
	router.HandleFunc("/", hc.get)

	ac := new(aboutController)
	ac.template = templates.Lookup("about.gohtml")
	ac.authorizer = authorizer
	router.HandleFunc("/about", ac.get)

	lc := new(loginController)
	lc.template = templates.Lookup("login.gohtml")
	lc.authorizer = authorizer
	router.HandleFunc("/login", lc.get).Methods("GET")
	router.HandleFunc("/login", lc.post).Methods("POST")

	loc := new(logoutController)
	loc.authorizer = authorizer
	router.HandleFunc("/logout", loc.get)

	sc := new(settingsController)
	sc.template = templates.Lookup("settings.gohtml")
	sc.authorizer = authorizer
	sc.DB = db
	router.HandleFunc("/settings", sc.get)

	http.Handle("/", router)

	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
	http.HandleFunc("/js/", serveResource)
}

func serveResource(w http.ResponseWriter, req *http.Request) {
	path := "public" + req.URL.Path
	var contentType string
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "text/javascript"
	} else {
		contentType = "text/plain"
	}

	f, err := os.Open(path)

	if err == nil {
		defer f.Close()
		w.Header().Add("Content-Type", contentType)

		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
	}
}
