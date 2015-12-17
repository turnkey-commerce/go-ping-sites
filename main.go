package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/turnkey-commerce/go-ping-sites/controllers"
	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/notifier"
	"github.com/turnkey-commerce/go-ping-sites/pinger"
)

func main() {
	var err error
	var db *sql.DB
	pinger.CreatePingerLog("")
	db, err = database.InitializeDB("go-ping-sites.db", "db-seed.toml")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()
	// Start the Pinger
	p := pinger.NewPinger(db, pinger.GetSites, pinger.RequestURL, pinger.DoExit,
		notifier.SendEmail, notifier.SendSms)
	p.Start()
	// Start the web server.
	templates := populateTemplates()
	controllers.Register(templates)
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
