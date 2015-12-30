package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"text/template"

	"golang.org/x/crypto/bcrypt"

	"github.com/apexskier/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/controllers"
	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/notifier"
	"github.com/turnkey-commerce/go-ping-sites/pinger"
)

var (
	authBackend     httpauth.SqlAuthBackend
	authBackendFile = "go-ping-sites-auth.db"
	roles           map[string]httpauth.Role
	authorizer      httpauth.Authorizer
)

func main() {
	var err error
	// Setup the main db.
	var db *sql.DB
	pinger.CreatePingerLog("")
	db, err = database.InitializeDB("go-ping-sites.db", "db-seed.toml")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()
	// Setup the auth
	// create the authorization backend
	authBackend, err = createAuthBackendFile()
	if err != nil {
		log.Fatal("Failed to create Auth Backend: ", err)
	}
	// create some default roles
	roles = make(map[string]httpauth.Role)
	roles["user"] = 30
	roles["admin"] = 80
	authorizer, err = httpauth.NewAuthorizer(authBackend, []byte("cookie-encryption-key"), "user", roles)
	createDefaultUser()
	// Start the Pinger
	p := pinger.NewPinger(db, pinger.GetSites, pinger.RequestURL, pinger.DoExit,
		notifier.SendEmail, notifier.SendSms)
	p.Start()
	// Start the web server.
	templates := populateTemplates()
	controllers.Register(db, templates)
	err = http.ListenAndServe(":8000", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func populateTemplates() *template.Template {
	result := template.New("templates")

	basePath := "templates"
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

	//fmt.Println(*templatePaths)
	result.ParseFiles(*templatePaths...)
	return result
}

func createDefaultUser() {
	// create a default user
	userName := "admin"
	if _, err := authBackend.User(userName); err != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte("adminpassword"), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		defaultUser := httpauth.UserData{Username: userName, Email: "admin@localhost.com", Hash: hash, Role: "admin"}
		err = authBackend.SaveUser(defaultUser)
		if err != nil {
			panic(err)
		}
	}
}

func createAuthBackendFile() (s httpauth.SqlAuthBackend, err error) {
	if _, err = os.Stat(authBackendFile); os.IsNotExist(err) {
		_, err = os.Create(authBackendFile)
		if err != nil {
			return s, err
		}
	}
	s, err = httpauth.NewSqlAuthBackend("sqlite3", authBackendFile)
	return s, err
}
