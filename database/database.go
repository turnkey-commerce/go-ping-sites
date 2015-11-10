package database

import (
	"database/sql"
	"time"
	// Import the sqlite3 package as blank.
	_ "github.com/mattn/go-sqlite3"
)

// Site is the website that will be monitored.
type Site struct {
	SiteID              int
	Name                string
	IsActive            bool
	URL                 string
	PingIntervalSeconds int
	TimeoutSeconds      int
	Contacts            []Contact
}

// Contact is one of the contacts for a particular site.
type Contact struct {
	Name         string
	EmailAddress string
	SmsNumber    string
	SmsActive    bool
	EmailActive  bool
}

// Ping contains information about a request to ping a site and details about the result
type Ping struct {
	SiteID         int
	TimeRequest    time.Time
	TimeResponse   time.Time
	HTTPStatusCode int
	TimedOut       bool
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

// GetSite gets a site and its collection of contacts from the DB by SiteID.
func (s *Site) GetSite(db *sql.DB, siteID int64) error {
	err := db.QueryRow("SELECT SiteID, Name, IsActive, URL, PingIntervalSeconds, TimeoutSeconds FROM Sites WHERE SiteID = $1", siteID).
		Scan(&s.SiteID, &s.Name, &s.IsActive, &s.URL, &s.PingIntervalSeconds, &s.TimeoutSeconds)
	if err != nil {
		return err
	}

	rows, err := db.Query(`SELECT Name, EmailAddress, SmsNumber, EmailActive, SmsActive
		FROM Contacts c JOIN  SiteContacts s  ON s.ContactID = c.ContactID WHERE s.siteID = $1`, siteID)

	for rows.Next() {
		var Name string
		var EmailAddress string
		var SmsNumber string
		var EmailActive bool
		var SmsActive bool
		err = rows.Scan(&Name, &EmailAddress, &SmsNumber, &EmailActive, &SmsActive)
		if err != nil {
			return err
		}
		s.Contacts = append(s.Contacts, Contact{Name: Name, EmailAddress: EmailAddress, SmsNumber: SmsNumber, EmailActive: EmailActive, SmsActive: SmsActive})
	}

	return nil
}

// CreateContact inserts a new contact in the DB and associates it with a site.
func (c Contact) CreateContact(db *sql.DB, siteID int64) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	result, err := db.Exec(
		"INSERT INTO Contacts (Name, EmailAddress, SmsNumber, EmailActive, SmsActive) VALUES ($1, $2, $3, $4, $5)",
		c.Name,
		c.EmailAddress,
		c.SmsNumber,
		c.EmailActive,
		c.SmsActive,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	contactID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Insert the contactID and the siteID in the many-to-many table
	result, errSiteContacts := db.Exec(
		"INSERT INTO SiteContacts (ContactID, SiteID) VALUES ($1, $2)",
		contactID,
		siteID,
	)
	if errSiteContacts != nil {
		tx.Rollback()
		return 0, errSiteContacts
	}

	tx.Commit()
	return contactID, nil
}

//CreatePing inserts a new ping row in the DB.
func (p Ping) CreatePing(db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO Pings (SiteID, TimeRequest, TimeResponse, HttpStatusCode, TimedOut) VALUES ($1, $2, $3, $4, $5)",
		p.SiteID,
		p.TimeRequest,
		p.TimeResponse,
		p.HTTPStatusCode,
		p.TimedOut,
	)
	if err != nil {
		return err
	}

	return nil
}
