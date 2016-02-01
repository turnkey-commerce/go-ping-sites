package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/apexskier/httpauth"
	"github.com/turnkey-commerce/go-ping-sites/config"
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

	roles = getRoles()
	authorizer, err = httpauth.NewAuthorizer(authBackend, []byte(config.Settings.Website.CookieKey), "user", roles)
	createDefaultUser()
	// Start the Pinger
	p := pinger.NewPinger(db, pinger.GetSites, pinger.RequestURL, pinger.DoExit,
		notifier.SendEmail, notifier.SendSms)
	p.Start()
	// Start the web server.
	templates := controllers.PopulateTemplates("templates")
	controllers.Register(db, authorizer, authBackend, roles, templates, p)
	err = http.ListenAndServe(":"+config.Settings.Website.HTTPPort, nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
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

func getRoles() map[string]httpauth.Role {
	// create some default roles
	var r = make(map[string]httpauth.Role)
	r["user"] = 30
	r["admin"] = 80
	return r
}
