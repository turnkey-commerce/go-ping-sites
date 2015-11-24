package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// The createStatements is used to initialize the DB with the schema.
// It seems better to put it in the code rather than an external file to
// prevent accidental changes by the users.
const createStatements = `CREATE TABLE "Sites" (
	"SiteId"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"Name"	TEXT NOT NULL UNIQUE,
	"IsActive"	INTEGER NOT NULL DEFAULT 1,
	"URL"	TEXT NOT NULL UNIQUE,
	"PingIntervalSeconds"	INTEGER NOT NULL DEFAULT 60,
	"TimeoutSeconds"	INTEGER NOT NULL DEFAULT 30
);

CREATE TABLE "SiteContacts" (
	"ContactId" INTEGER NOT NULL,
	"SiteId" INTEGER NOT NULL,
	FOREIGN KEY("ContactId")	REFERENCES "Contacts"("ContactId"),
	FOREIGN KEY("SiteId")	REFERENCES "Sites"("SiteId")
	PRIMARY KEY("ContactId","SiteId")
);

CREATE TABLE "Pings" (
	"TimeRequest"	TIMESTAMP NOT NULL,
	"SiteId"	INTEGER NOT NULL,
	"TimeResponse"	TIMESTAMP,
	"HttpStatusCode"	INTEGER,
	"TimedOut"	INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("TimeRequest","SiteId")
	FOREIGN KEY("SiteId") REFERENCES "Sites"("SiteId")
);

CREATE TABLE "Contacts" (
	"ContactId"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"Name"	TEXT NOT NULL UNIQUE,
	"EmailAddress"	TEXT NOT NULL,
	"SmsNumber"	INTEGER,
	"SmsActive"	INTEGER NOT NULL DEFAULT 0,
	"EmailActive"	INTEGER NOT NULL DEFAULT 1
);`

// InitializeDB creates the DB file and the schema if the file doesn't exist.
func InitializeDB(dbPath string, seedFile string) (*sql.DB, error) {
	newDB := false
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		newDB = true
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if newDB {
		fmt.Println("New Database, creating Schema...")
		err = createSchema(db)
		if err != nil {
			return nil, err
		}
		// If a seed config file exists then use it to seed the initial DB.
		if _, err := os.Stat(seedFile); err == nil {
			err = seedInitialSites(db, seedFile)
			if err != nil {
				return nil, err
			}
		}
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	return db, nil
}

// createSchema applies the initial schema creation to the database.
func createSchema(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = db.Exec(createStatements)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

//Seed represents the initial seed to the DB.~
var Seed struct {
	Sites    []Site
	Contacts []Contact
}

// seedInitialSites gets some initial sites from a config file
func seedInitialSites(db *sql.DB, seedFile string) error {
	fmt.Println("Seeding initial sites with", seedFile, "...")
	var err error
	if _, err = toml.DecodeFile(seedFile, &Seed); err != nil {
		log.Println(err)
		return err
	}

	for i, s := range Seed.Sites {
		err = s.CreateSite(db)
		Seed.Sites[i].SiteID = s.SiteID
		if err != nil {
			return err
		}
	}
	for _, c := range Seed.Contacts {
		err = c.CreateContact(db)
		if err != nil {
			return err
		}
		for _, s := range Seed.Sites {
			err = c.AddContactToSite(db, s.SiteID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

const testDb string = "./test.db"

// InitializeTestDB is for test packages to initalize a DB for integration testing.
func InitializeTestDB(seedFile string) (*sql.DB, error) {
	var db *sql.DB
	err := deleteDb(testDb)
	if err != nil {
		return nil, err
	}
	db, err = InitializeDB(testDb, seedFile)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// deleteDb removes the DB file, mainly intended for testing
func deleteDb(dbPath string) error {
	if _, err := os.Stat(dbPath); err == nil {
		err := os.Remove(dbPath)
		if err != nil {
			return err
		}
	}
	return nil
}
