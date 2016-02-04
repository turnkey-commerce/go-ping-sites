package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/turnkey-commerce/go-ping-sites/config"
	"github.com/turnkey-commerce/go-ping-sites/controllers"
	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/notifier"
	"github.com/turnkey-commerce/go-ping-sites/pinger"
	"github.com/turnkey-commerce/httpauth"
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
		fatalError("Failed to initialize database:", err)
	}
	defer db.Close()
	// Validate the required config for the website.
	_, err = govalidator.ValidateStruct(config.Settings.Website)
	if err != nil {
		fatalError("Error in config.toml configuration:", err)
	}
	// Setup the auth
	// create the authorization backend
	authBackend, err = createAuthBackendFile()
	if err != nil {
		fatalError("Failed to create Auth Backend: ", err)
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
		fatalError("ListenAndServe: ", err)
	}
}

func createDefaultUser() {
	// create a default user
	userName := "admin"
	pwd := "adminpassword"
	if _, err := authBackend.User(userName); err != nil {
		defaultUser := httpauth.UserData{Username: userName, Email: "admin@localhost.com", Role: "admin"}
		err = authBackend.SaveUser(defaultUser)
		if err != nil {
			panic(err)
		}
		err = authorizer.Update(nil, nil, userName, pwd, "")
		if err != nil {
			panic(err)
		}
	}
}

func fatalError(message string, content interface{}) {
	fmt.Println(message, content)
	log.Fatal(message, content)
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
