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
	"github.com/turnkey-commerce/go-ping-sites/pinger"
)

// CurrentUserGetter gets the current user from the http request
type CurrentUserGetter interface {
	CurrentUser(rw http.ResponseWriter, req *http.Request) (user httpauth.UserData, e error)
	Messages(rw http.ResponseWriter, req *http.Request) []string
}

// Register the handlers for a given route.
func Register(db *sql.DB, authorizer httpauth.Authorizer, authBackend httpauth.AuthBackend,
	roles map[string]httpauth.Role, templates *template.Template, pinger *pinger.Pinger,
	version string) {
	router := mux.NewRouter()

	hc := new(homeController)
	hc.template = templates.Lookup("home.gohtml")
	hc.authorizer = authorizer
	hc.DB = db
	router.Handle("/", authorizeRole(http.HandlerFunc(hc.get), authorizer, "user"))

	rc := new(reportsController)
	rc.template = templates.Lookup("reports.gohtml")
	rc.authorizer = authorizer
	rc.DB = db
	router.Handle("/reports", authorizeRole(appHandler(rc.get), authorizer, "user"))

	ac := new(aboutController)
	ac.template = templates.Lookup("about.gohtml")
	ac.authorizer = authorizer
	ac.version = version
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
	router.Handle("/settings", authorizeRole(appHandler(sc.get), authorizer, "admin"))

	//settingsSub is a subrouter "/settings"
	settingsSub := router.PathPrefix("/settings").Subrouter()

	// /settings/users
	uc := new(usersController)
	uc.getTemplate = templates.Lookup("users.gohtml")
	uc.editTemplate = templates.Lookup("user_edit.gohtml")
	uc.newTemplate = templates.Lookup("user_new.gohtml")
	uc.authorizer = authorizer
	uc.authBackend = authBackend
	uc.roles = roles
	settingsSub.Handle("/users", authorizeRole(appHandler(uc.get), authorizer, "admin"))
	settingsSub.Handle("/users/{username}/edit", authorizeRole(appHandler(uc.editGet), authorizer, "admin")).Methods("GET")
	settingsSub.Handle("/users/{username}/edit", authorizeRole(appHandler(uc.editPost), authorizer, "admin")).Methods("POST")
	settingsSub.Handle("/users/new", authorizeRole(appHandler(uc.newGet), authorizer, "admin")).Methods("GET")
	settingsSub.Handle("/users/new", authorizeRole(appHandler(uc.newPost), authorizer, "admin")).Methods("POST")

	// /settings/contacts
	cc := new(contactsController)
	cc.getTemplate = templates.Lookup("contacts.gohtml")
	cc.editTemplate = templates.Lookup("contact_edit.gohtml")
	cc.newTemplate = templates.Lookup("contact_new.gohtml")
	cc.authorizer = authorizer
	cc.pinger = pinger
	cc.DB = db
	settingsSub.Handle("/contacts", authorizeRole(appHandler(cc.get), authorizer, "admin"))
	settingsSub.Handle("/contacts/{contactID}/edit", authorizeRole(appHandler(cc.editGet), authorizer, "admin")).Methods("GET")
	settingsSub.Handle("/contacts/{contactID}/edit", authorizeRole(appHandler(cc.editPost), authorizer, "admin")).Methods("POST")
	settingsSub.Handle("/contacts/new", authorizeRole(appHandler(cc.newGet), authorizer, "admin")).Methods("GET")
	settingsSub.Handle("/contacts/new", authorizeRole(appHandler(cc.newPost), authorizer, "admin")).Methods("POST")

	// /settings/sites
	stc := new(sitesController)
	stc.detailsTemplate = templates.Lookup("site_details.gohtml")
	stc.editTemplate = templates.Lookup("site_edit.gohtml")
	stc.newTemplate = templates.Lookup("site_new.gohtml")
	stc.changeContactsTemplate = templates.Lookup("site_change_contacts.gohtml")
	stc.authorizer = authorizer
	stc.pinger = pinger
	stc.DB = db
	// Site list is handled on main settings page so not required here.
	settingsSub.Handle("/sites/new", authorizeRole(appHandler(stc.newGet), authorizer, "admin")).Methods("GET")
	settingsSub.Handle("/sites/new", authorizeRole(appHandler(stc.newPost), authorizer, "admin")).Methods("POST")
	settingsSub.Handle("/sites/{siteID}", authorizeRole(appHandler(stc.getDetails), authorizer, "admin"))
	settingsSub.Handle("/sites/{siteID}/edit", authorizeRole(appHandler(stc.editGet), authorizer, "admin")).Methods("GET")
	settingsSub.Handle("/sites/{siteID}/edit", authorizeRole(appHandler(stc.editPost), authorizer, "admin")).Methods("POST")

	http.Handle("/", router)

	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
	http.HandleFunc("/js/", serveResource)
	http.HandleFunc("/fonts/", serveResource)
}

func serveResource(w http.ResponseWriter, req *http.Request) {
	path := "public" + req.URL.Path
	var contentType string
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(path, ".eot") {
		contentType = "application/vnd.ms-fontobject"
	} else if strings.HasSuffix(path, ".ttf") {
		contentType = "application/font-sfnt"
	} else if strings.HasSuffix(path, ".otf") {
		contentType = "application/font-sfnt"
	} else if strings.HasSuffix(path, ".woff") {
		contentType = "application/font-woff"
	} else if strings.HasSuffix(path, ".woff2") {
		contentType = "application/font-woff2"
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

func authorizeRole(h http.Handler, authorizer httpauth.Authorizer, role string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if err := authorizer.AuthorizeRole(rw, req, role, true); err != nil {
			// Redirect to about in  to avoid confusing the user if it's about privileges
			// This also avoids a redirect loop if the main dashboard page is not authorized.
			if strings.Contains(err.Error(), "user not logged in") ||
				strings.Contains(err.Error(), "new authorization session") {
				http.Redirect(rw, req, "/login", http.StatusSeeOther)
			} else {
				http.Redirect(rw, req, "/about", http.StatusSeeOther)
			}
			return
		}
		h.ServeHTTP(rw, req)
	})
}

type appHandler func(http.ResponseWriter, *http.Request) (int, error)

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := fn(w, r); err != nil {
		switch status {
		// We can have cases as granular as we like, if we wanted to
		// return custom errors for specific status codes.
		case http.StatusInternalServerError:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			// Catch any other errors we haven't explicitly handled
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func getCurrentUser(rw http.ResponseWriter, req *http.Request, authorizer CurrentUserGetter) (isAuthenticated bool, user httpauth.UserData) {
	isAuthenticated = false
	var err error
	user, err = authorizer.CurrentUser(rw, req)
	if err == nil {
		isAuthenticated = true
	}
	return isAuthenticated, user
}

// PopulateTemplates loads and parses all of the templates in the templates directory
func PopulateTemplates(templatePath string) *template.Template {
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

	result.ParseFiles(*templatePaths...)
	return result
}
