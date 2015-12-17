package controllers

import (
	"bufio"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

// Register the handlers for a given route.
func Register(templates *template.Template) {
	router := mux.NewRouter()

	hc := new(homeController)
	hc.template = templates.Lookup("home.html")
	router.HandleFunc("/home", hc.get)

	ac := new(aboutController)
	ac.template = templates.Lookup("about.html")
	router.HandleFunc("/about", ac.get)

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
