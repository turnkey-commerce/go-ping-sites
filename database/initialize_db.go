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
	"SiteId"	           INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"Name"	               TEXT NOT NULL UNIQUE,
	"IsActive"	           INTEGER NOT NULL DEFAULT 1,
	"URL"	               TEXT NOT NULL UNIQUE,
	"PingIntervalSeconds"  INTEGER NOT NULL DEFAULT 60,
	"TimeoutSeconds"	   INTEGER NOT NULL DEFAULT 30,
	"IsSiteUp"             INTEGER NOT NULL DEFAULT 1,
	"LastStatusChange"     TIMESTAMP,
	"LastPing"             TIMESTAMP
);

CREATE TABLE "SiteContacts" (
	"ContactId"              INTEGER NOT NULL,
	"SiteId"                 INTEGER NOT NULL,
	FOREIGN KEY("ContactId") REFERENCES "Contacts"("ContactId"),
	FOREIGN KEY("SiteId")	 REFERENCES "Sites"("SiteId")
	PRIMARY KEY("ContactId","SiteId")
);

CREATE TABLE "Pings" (
	"TimeRequest"	  TIMESTAMP NOT NULL,
	"SiteId"	      INTEGER NOT NULL,
	"Duration"        INTEGER,
	"HttpStatusCode"  INTEGER,
	"TimedOut"	      INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("TimeRequest","SiteId")
	FOREIGN KEY("SiteId") REFERENCES "Sites"("SiteId")
);

CREATE TABLE "Contacts" (
	"ContactId"	    INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"Name"	        TEXT NOT NULL UNIQUE,
	"EmailAddress"  TEXT NOT NULL,
	"SmsNumber"	    INTEGER,
	"SmsActive"	    INTEGER NOT NULL DEFAULT 0,
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
	}

	err = upgradeDB(db)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	// Do the seeding, if applicable, after the upgrade to avoid schema issues.
	if newDB {
		// If a seed config file exists then use it to seed the initial DB.
		if _, err := os.Stat(seedFile); err == nil {
			err = seedInitialSites(db, seedFile)
			if err != nil {
				return nil, err
			}
		}
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
			err = s.AddContactToSite(db, c.ContactID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// upgradeStatements is used to upgrade the DB from the initial state that
// was created in the createStatements. If new statements are added then then
// databaseVersion constant should be incremented below by 1.
const upgradeStatements1 = `
	CREATE INDEX IF NOT EXISTS pings_timerequest_httpstatuscode
	ON pings (TimeRequest, HttpStatusCode);
	ALTER TABLE "Sites" ADD COLUMN "FirstPing" TIMESTAMP;
	UPDATE Sites SET FirstPing = '0001-01-01 00:00:00+00:00' WHERE FirstPing IS NULL;
`

const upgradeStatements2 = `
	ALTER TABLE "Sites" ADD COLUMN "ContentExpected"    TEXT NOT NULL DEFAULT '';
	ALTER TABLE "Sites" ADD COLUMN "ContentNotExpected" TEXT NOT NULL DEFAULT '';
`

// If new upgrade statements are added then this must be incremented by 1.
const databaseVersion int32 = 3

//upgradeDB applies any upgrades since the initial schema of the DB.
func upgradeDB(db *sql.DB) error {
	// first check if upgrade is necessary
	var currentVersion int32
	err := db.QueryRow("PRAGMA user_version;").Scan(&currentVersion)
	if err != nil {
		return err
	}

	if currentVersion == databaseVersion {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if currentVersion < 2 {
		_, err = db.Exec(upgradeStatements1)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if currentVersion < 3 {
		_, err = db.Exec(upgradeStatements2)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = db.Exec(fmt.Sprintf("PRAGMA user_version = %d", databaseVersion))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	log.Println("Upgraded database to version ", databaseVersion)
	return nil
}

const testDb string = "./test.db"

// InitializeTestDB is for test packages to initialize a DB for integration testing.
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

const reportDB string = "./report_test.db"

// InitializeReportDB is for test packages to initialize a Report DB for integration testing.
func InitializeReportDB() (*sql.DB, error) {
	var db *sql.DB
	db, err := InitializeDB(reportDB, "")
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
