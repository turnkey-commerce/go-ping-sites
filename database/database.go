package database

import (
	"database/sql"
	"time"
	// Import the sqlite3 package as blank.
	_ "github.com/mattn/go-sqlite3"
)

// Site is the website that will be monitored.
type Site struct {
	SiteID              int64
	Name                string
	IsActive            bool
	URL                 string
	PingIntervalSeconds int
	TimeoutSeconds      int
	Contacts            []Contact
	Pings               []Ping
}

// Contact is one of the contacts for a particular site.
type Contact struct {
	ContactID    int64
	Name         string
	EmailAddress string
	SmsNumber    string
	SmsActive    bool
	EmailActive  bool
}

// Sites is a slice of sites
type Sites []Site

// Contacts is a slice of contacts that aren't necessarily associated with a given site.
type Contacts []Contact

// Ping contains information about a request to ping a site and details about the result
type Ping struct {
	SiteID         int64
	TimeRequest    time.Time
	Duration       int
	HTTPStatusCode int
	TimedOut       bool
}

//CreateSite inserts a new site in the DB.
func (s *Site) CreateSite(db *sql.DB) error {
	result, err := db.Exec(
		`INSERT INTO Sites (Name, IsActive, URL, PingIntervalSeconds, TimeoutSeconds)
			VALUES ($1, $2, $3, $4, $5)`,
		s.Name,
		s.IsActive,
		s.URL,
		s.PingIntervalSeconds,
		s.TimeoutSeconds,
	)
	if err != nil {
		return err
	}

	s.SiteID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

// GetSite gets the site details for a given site.
func (s *Site) GetSite(db *sql.DB, siteID int64) error {
	err := db.QueryRow(`SELECT SiteID, Name, IsActive, URL, PingIntervalSeconds,
		TimeoutSeconds FROM Sites WHERE SiteID = $1`, siteID).
		Scan(&s.SiteID, &s.Name, &s.IsActive, &s.URL, &s.PingIntervalSeconds, &s.TimeoutSeconds)
	if err != nil {
		return err
	}
	return nil
}

// GetActiveSitesWithContacts gets all of the active sites with the contacts.
func (s *Sites) GetActiveSitesWithContacts(db *sql.DB) error {
	rows, err := db.Query(`SELECT SiteID, Name, IsActive, URL, PingIntervalSeconds,
		TimeoutSeconds FROM Sites WHERE IsActive = $1
		ORDER BY Name`, true)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var SiteID int64
		var Name string
		var IsActive bool
		var URL string
		var PingIntervalSeconds int
		var TimeoutSeconds int
		err = rows.Scan(&SiteID, &Name, &IsActive, &URL, &PingIntervalSeconds, &TimeoutSeconds)
		if err != nil {
			return err
		}
		site := Site{SiteID: SiteID, Name: Name, IsActive: IsActive, URL: URL,
			PingIntervalSeconds: PingIntervalSeconds, TimeoutSeconds: TimeoutSeconds}
		err = site.GetSiteContacts(db, site.SiteID)
		if err != nil {
			return err
		}
		*s = append(*s, site)
	}
	return nil
}

// GetSiteContacts gets the collection of contacts for a given site.
func (s *Site) GetSiteContacts(db *sql.DB, siteID int64) error {
	rows, err := db.Query(`SELECT c.ContactID, Name, EmailAddress, SmsNumber, EmailActive, SmsActive
		FROM Contacts c JOIN  SiteContacts s  ON s.ContactID = c.ContactID WHERE s.siteID = $1
		ORDER BY Name`, siteID)
	if err != nil {
		return err
	}

	// nil out the slice in case it is rereading it from the DB.
	s.Contacts = nil
	defer rows.Close()
	for rows.Next() {
		var ContactID int64
		var Name string
		var EmailAddress string
		var SmsNumber string
		var EmailActive bool
		var SmsActive bool
		err = rows.Scan(&ContactID, &Name, &EmailAddress, &SmsNumber, &EmailActive, &SmsActive)
		if err != nil {
			return err
		}
		s.Contacts = append(s.Contacts, Contact{ContactID: ContactID, Name: Name,
			EmailAddress: EmailAddress, SmsNumber: SmsNumber, EmailActive: EmailActive,
			SmsActive: SmsActive})
	}

	return nil
}

// CreateContact inserts a new contact in the DB.
func (c *Contact) CreateContact(db *sql.DB) error {
	result, err := db.Exec(
		"INSERT INTO Contacts (Name, EmailAddress, SmsNumber, EmailActive, SmsActive) VALUES ($1, $2, $3, $4, $5)",
		c.Name,
		c.EmailAddress,
		c.SmsNumber,
		c.EmailActive,
		c.SmsActive,
	)
	if err != nil {
		return err
	}

	c.ContactID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

// AddContactToSite associates a contact with a site.
func (c Contact) AddContactToSite(db *sql.DB, siteID int64) error {
	// Insert the contactID and the siteID in the many-to-many table
	_, err := db.Exec(
		"INSERT INTO SiteContacts (ContactID, SiteID) VALUES ($1, $2)",
		c.ContactID,
		siteID,
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveContactFromSite deletes the association of a contact with a site.
func (c Contact) RemoveContactFromSite(db *sql.DB, siteID int64) error {
	// Insert the contactID and the siteID in the many-to-many table
	_, err := db.Exec(
		"DELETE FROM SiteContacts WHERE ContactID = $1 AND SiteID = $2",
		c.ContactID,
		siteID,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetContacts gets all contacts
func (c *Contacts) GetContacts(db *sql.DB) error {
	rows, err := db.Query(`SELECT ContactID, Name, EmailAddress, SmsNumber, EmailActive, SmsActive
		FROM Contacts
	  ORDER BY Name`)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var ContactID int64
		var Name string
		var EmailAddress string
		var SmsNumber string
		var EmailActive bool
		var SmsActive bool
		err = rows.Scan(&ContactID, &Name, &EmailAddress, &SmsNumber, &EmailActive, &SmsActive)
		if err != nil {
			return err
		}
		*c = append(*c, Contact{ContactID: ContactID, Name: Name,
			EmailAddress: EmailAddress, SmsNumber: SmsNumber, EmailActive: EmailActive,
			SmsActive: SmsActive})
	}

	return nil
}

//CreatePing inserts a new ping row in the DB.
func (p Ping) CreatePing(db *sql.DB) error {
	_, err := db.Exec(
		`INSERT INTO Pings (SiteID, TimeRequest, Duration, HttpStatusCode, TimedOut)
			VALUES ($1, $2, $3, $4, $5)`,
		p.SiteID,
		p.TimeRequest,
		p.Duration,
		p.HTTPStatusCode,
		p.TimedOut,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetSitePings gets the pings for a given site for a given time interval.
func (s *Site) GetSitePings(db *sql.DB, siteID int64, startTime time.Time, endTime time.Time) error {
	rows, err := db.Query(`SELECT SiteID, TimeRequest, Duration, HttpStatusCode, TimedOut
		FROM Pings WHERE SiteID = $1 AND TimeRequest >= $2 AND TimeRequest <=$3
		ORDER BY TimeRequest`, siteID, startTime, endTime)
	if err != nil {
		return err
	}

	// nil out the slice in case it is rereading it from the DB.
	s.Pings = nil
	defer rows.Close()
	for rows.Next() {
		var SiteID int64
		var TimeRequest time.Time
		var Duration int
		var HTTPStatusCode int
		var TimedOut bool
		err = rows.Scan(&SiteID, &TimeRequest, &Duration, &HTTPStatusCode, &TimedOut)
		if err != nil {
			return err
		}
		s.Pings = append(s.Pings, Ping{SiteID: SiteID, TimeRequest: TimeRequest,
			Duration: Duration, HTTPStatusCode: HTTPStatusCode, TimedOut: TimedOut})
	}

	return nil
}
