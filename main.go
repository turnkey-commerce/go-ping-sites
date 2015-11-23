package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/notifier"
	"github.com/turnkey-commerce/go-ping-sites/pinger"
)

func main() {
	var err error
	var db *sql.DB
	pinger.CreatePingerLog("")
	db, err = database.InitializeDB("go-ping-sites.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()
	p := pinger.NewPinger(db, pinger.GetSites, pinger.RequestURL, pinger.DoExit,
		notifier.SendEmail, notifier.SendSms)
	p.Start()
	fmt.Scanln()
}
