package database

import (
	"database/sql"
	"fmt"
	"os"
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
	"Name"	TEXT NOT NULL,
	"EmailAddress"	TEXT NOT NULL,
	"SmsNumber"	INTEGER,
	"SmsActive"	INTEGER NOT NULL DEFAULT 0,
	"EmailActive"	INTEGER NOT NULL DEFAULT 1
);`

// InitializeDB creates the DB file and the schema if the file doesn't exist.
func InitializeDB(dbPath string) (*sql.DB, error) {
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
		err = CreateSchema(db)
		if err != nil {
			return nil, err
		}
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateSchema applies the initial schema creation to the database.
func CreateSchema(db *sql.DB) error {
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

// DeleteDb removes the DB file, mainly intended for testing
func DeleteDb(dbPath string) error {
	if _, err := os.Stat(dbPath); err == nil {
		err := os.Remove(dbPath)
		if err != nil {
			return err
		}
	}
	return nil
}
