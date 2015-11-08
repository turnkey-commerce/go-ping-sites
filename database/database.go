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
	dat, err := ioutil.ReadFile(dbCreatePath)
	if err != nil {
		return nil, err
	}
	sqlCreate := (string(dat))

	_, err = db.Exec(sqlCreate)
	if err != nil {
		return nil, err
	}
	return db, nil
}
