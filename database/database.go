package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	// Import the sqlite3 package as blank.
	_ "github.com/mattn/go-sqlite3"
)

//Site is the website that will be monitored.
type Site struct {
	Name                string
	IsActive            bool
	URL                 string
	PingIntervalSeconds int
	TimeoutSeconds      int
}

// InitializeDB creates the DB file and imports the schema.
func InitializeDB() (*sql.DB, error) {
	dbPath := "./go-ping-sites.db"
	os.Remove(dbPath)
	fmt.Println("Removed file")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	dbCreatePath, _ := filepath.Abs("../database/create_database.sql")
	createStatements, err := ioutil.ReadFile(dbCreatePath)
	if err != nil {
		return nil, err
	}
	sqlCreate := (string(createStatements))

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(sqlCreate)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return db, nil
}

//CreateSite inserts a new site in the DB.
func (s Site) CreateSite(db *sql.DB) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO Sites (Name, IsActive, URL, PingIntervalSeconds, TimeoutSeconds) VALUES ($1, $2, $3, $4, $5)",
		s.Name,
		s.IsActive,
		s.URL,
		s.PingIntervalSeconds,
		s.TimeoutSeconds,
	)
	if err != nil {
		return 0, err
	}

	siteID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return siteID, nil
}

//GetSite gets a site from the DB by SiteID.
func (s *Site) GetSite(db *sql.DB, siteID int64) error {
	err := db.QueryRow("SELECT Name, IsActive, URL, PingIntervalSeconds, TimeoutSeconds FROM Sites WHERE SiteID = $1", siteID).
		Scan(&s.Name, &s.IsActive, &s.URL, &s.PingIntervalSeconds, &s.TimeoutSeconds)
	if err != nil {
		return err
	}
	return nil
}
