package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/turnkey-commerce/go-ping-sites/database"
	"github.com/turnkey-commerce/go-ping-sites/pinger"
)

func main() {
	var err error
	var db *sql.DB
	db, err = database.InitializeDB("go-ping-sites.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()
	p := pinger.NewPinger(db, getSites, requestURL)
	p.Start()
	fmt.Scanln()
}

func requestURL(url string, timeout int) (string, int, error) {
	to := time.Duration(timeout) * time.Second
	client := http.Client{
		Timeout: to,
	}
	res, err := client.Get(url)
	if err != nil {
		return "", 0, err
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", 0, err
	}
	return string(content), res.StatusCode, nil
}

func getSites(db *sql.DB) (database.Sites, error) {
	var sites database.Sites
	err := sites.GetActiveSitesWithContacts(db)
	if err != nil {
		return nil, err
	}
	return sites, nil
}
