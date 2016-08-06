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
	IsSiteUp            bool
	ContentExpected     string
	ContentUnexpected   string
	LastStatusChange    time.Time
	LastPing            time.Time
	FirstPing           time.Time
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
	SiteCount    int
	Sites        []Site
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

// Report contains information about performance where AvgResponse is the average
// response time for successful requests, PingsUp are the number of successful
// pings when the site was up and PingsDown is the number of pings when the site
// was down.
type Report struct {
	AvgResponse float64
	PingsUp     int
	PingsDown   int
}

// Reports is a slice of reports, usually the index will represent the month.
type Reports []Report

//CreateSite inserts a new site in the DB.
func (s *Site) CreateSite(db *sql.DB) error {
	// Set site to initially be up, as is the assumption when the pinging first starts.
	s.IsSiteUp = true
	result, err := db.Exec(
		`INSERT INTO Sites (Name, IsActive, URL, PingIntervalSeconds, TimeoutSeconds,
			IsSiteUp, LastStatusChange, LastPing, FirstPing)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		s.Name,
		s.IsActive,
		s.URL,
		s.PingIntervalSeconds,
		s.TimeoutSeconds,
		s.IsSiteUp,
		s.LastStatusChange,
		s.LastPing,
		s.FirstPing,
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

// UpdateSite updates the site information in the DB.
func (s *Site) UpdateSite(db *sql.DB) error {
	_, err := db.Exec(
		`Update Sites SET Name = $1, URL = $2, IsActive = $3,
		  PingIntervalSeconds = $4, TimeoutSeconds = $5
			WHERE SiteId = $6`,
		s.Name,
		s.URL,
		s.IsActive,
		s.PingIntervalSeconds,
		s.TimeoutSeconds,
		s.SiteID,
	)
	if err != nil {
		return err
	}
	return nil
}

//UpdateSiteStatus updates the up/down status and last status change of a Site.
func (s *Site) UpdateSiteStatus(db *sql.DB, isSiteUp bool) error {
	_, err := db.Exec(
		`UPDATE Sites SET IsSiteUp = $1, LastStatusChange = $2
			WHERE SiteId = $3`,
		isSiteUp,
		time.Now(),
		s.SiteID,
	)
	if err != nil {
		return err
	}

	return nil
}

//UpdateSiteFirstPing updates the up/down status and last status change of a Site.
func (s *Site) UpdateSiteFirstPing(db *sql.DB, firstPingTime time.Time) error {
	_, err := db.Exec(
		`UPDATE Sites SET FirstPing = $1
			WHERE SiteId = $2`,
		firstPingTime,
		s.SiteID,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetSite gets the site details for a given site.
func (s *Site) GetSite(db *sql.DB, siteID int64) error {
	err := db.QueryRow(`SELECT SiteID, Name, IsActive, URL, PingIntervalSeconds,
		TimeoutSeconds, IsSiteUp, LastStatusChange, LastPing, FirstPing FROM Sites
		WHERE SiteID = $1`, siteID).
		Scan(&s.SiteID, &s.Name, &s.IsActive, &s.URL, &s.PingIntervalSeconds, &s.TimeoutSeconds,
			&s.IsSiteUp, &s.LastStatusChange, &s.LastPing, &s.FirstPing)
	if err != nil {
		return err
	}
	return nil
}

const getActiveSitesQueryString string = `SELECT SiteID, Name, IsActive, URL,
	PingIntervalSeconds, TimeoutSeconds, IsSiteUp, LastStatusChange, LastPing,
	FirstPing FROM Sites WHERE IsActive = $1
	ORDER BY Name`

const getAllSitesQueryString string = `SELECT SiteID, Name, IsActive, URL,
  PingIntervalSeconds, TimeoutSeconds, IsSiteUp, LastStatusChange, LastPing,
	FirstPing FROM Sites
	ORDER BY Name`

// GetSites gets all of the sites without contacts.
// The switch activeOnly controls whether to get only active sites.
// The option withContacts controls whether to also get the associated contacts.
func (s *Sites) GetSites(db *sql.DB, activeOnly bool, withContacts bool) error {
	var queryString string
	if activeOnly {
		queryString = getActiveSitesQueryString
	} else {
		queryString = getAllSitesQueryString
	}
	rows, err := db.Query(queryString, true)
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
		var IsSiteUp bool
		var LastStatusChange time.Time
		var LastPing time.Time
		var FirstPing time.Time
		err = rows.Scan(&SiteID, &Name, &IsActive, &URL, &PingIntervalSeconds, &TimeoutSeconds,
			&IsSiteUp, &LastStatusChange, &LastPing, &FirstPing)
		if err != nil {
			return err
		}
		site := Site{SiteID: SiteID, Name: Name, IsActive: IsActive, URL: URL,
			PingIntervalSeconds: PingIntervalSeconds, TimeoutSeconds: TimeoutSeconds,
			IsSiteUp: IsSiteUp, LastStatusChange: LastStatusChange, LastPing: LastPing,
			FirstPing: FirstPing}
		if withContacts {
			err = site.GetSiteContacts(db, site.SiteID)
			if err != nil {
				return err
			}
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

// GetContactSites gets the collection of sites for a given contact.
func (c *Contact) GetContactSites(db *sql.DB) error {
	rows, err := db.Query(`SELECT s.SiteID, Name, IsActive, URL, PingIntervalSeconds,
		TimeoutSeconds, IsSiteUp, LastStatusChange, LastPing, FirstPing
		FROM Sites s JOIN SiteContacts sc ON s.SiteID = sc.SiteID WHERE sc.ContactID = $1
		ORDER BY Name`, c.ContactID)
	if err != nil {
		return err
	}
	// nil out the slice in case it is rereading it from the DB.
	c.Sites = nil
	defer rows.Close()
	for rows.Next() {
		var SiteID int64
		var Name string
		var IsActive bool
		var URL string
		var PingIntervalSeconds int
		var TimeoutSeconds int
		var IsSiteUp bool
		var LastStatusChange time.Time
		var LastPing time.Time
		var FirstPing time.Time
		err = rows.Scan(&SiteID, &Name, &IsActive, &URL, &PingIntervalSeconds,
			&TimeoutSeconds, &IsSiteUp, &LastStatusChange, &LastPing, &FirstPing)
		if err != nil {
			return err
		}
		c.Sites = append(c.Sites, Site{SiteID: SiteID, Name: Name, IsActive: IsActive, URL: URL,
			PingIntervalSeconds: PingIntervalSeconds, TimeoutSeconds: TimeoutSeconds,
			IsSiteUp: IsSiteUp, LastStatusChange: LastStatusChange, LastPing: LastPing,
			FirstPing: FirstPing})
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

// UpdateContact updates the contact information in the DB.
func (c *Contact) UpdateContact(db *sql.DB) error {
	_, err := db.Exec(
		`Update Contacts SET Name = $1, EmailAddress = $2, SmsNumber = $3,
		  EmailActive = $4, SmsActive = $5
			WHERE ContactID = $6`,
		c.Name,
		c.EmailAddress,
		c.SmsNumber,
		c.EmailActive,
		c.SmsActive,
		c.ContactID,
	)
	if err != nil {
		return err
	}
	return nil
}

// DeleteContact deletes the contact from the DB.
func (c *Contact) DeleteContact(db *sql.DB) error {
	// Do in a transaction because we have to first delete the contact on the sites.
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = db.Exec(
		`DELETE FROM SiteContacts WHERE ContactID = $1`,
		c.ContactID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = db.Exec(
		`DELETE FROM Contacts WHERE ContactID = $1;`,
		c.ContactID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// GetContact gets the contact details for a given contact.
func (c *Contact) GetContact(db *sql.DB, contactID int64) error {
	err := db.QueryRow(`SELECT ContactID, Name, EmailAddress, SmsNumber, SmsActive,
		EmailActive FROM Contacts WHERE ContactID = $1`, contactID).
		Scan(&c.ContactID, &c.Name, &c.EmailAddress, &c.SmsNumber, &c.SmsActive,
			&c.EmailActive)
	if err != nil {
		return err
	}
	return nil
}

// AddContactToSite associates a contact with a site.
func (s Site) AddContactToSite(db *sql.DB, contactID int64) error {
	// Insert the contactID and the siteID in the many-to-many table
	_, err := db.Exec(
		"INSERT INTO SiteContacts (ContactID, SiteID) VALUES ($1, $2)",
		contactID,
		s.SiteID,
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveContactFromSite deletes the association of a contact with a site.
func (s Site) RemoveContactFromSite(db *sql.DB, contactID int64) error {
	// Insert the contactID and the siteID in the many-to-many table
	_, err := db.Exec(
		"DELETE FROM SiteContacts WHERE ContactID = $1 AND SiteID = $2",
		contactID,
		s.SiteID,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetContacts gets all contacts
func (c *Contacts) GetContacts(db *sql.DB) error {
	rows, err := db.Query(`SELECT Contacts.ContactID, Name, EmailAddress, SmsNumber,
		EmailActive, SmsActive, count(Distinct SiteContacts.SiteID) AS SiteCount
		FROM Contacts LEFT JOIN SiteContacts ON Contacts.ContactId = SiteContacts.ContactId
		GROUP BY Contacts.ContactId, Name, EmailAddress, SmsNumber, EmailActive, SmsActive
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
		var SiteCount int
		err = rows.Scan(&ContactID, &Name, &EmailAddress, &SmsNumber, &EmailActive,
			&SmsActive, &SiteCount)
		if err != nil {
			return err
		}
		*c = append(*c, Contact{ContactID: ContactID, Name: Name,
			EmailAddress: EmailAddress, SmsNumber: SmsNumber, EmailActive: EmailActive,
			SmsActive: SmsActive, SiteCount: SiteCount})
	}

	return nil
}

//CreatePing inserts a new ping row and last check for the site in the DB.
func (p Ping) CreatePing(db *sql.DB) error {
	var err error
	_, err = db.Exec(
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

	_, err = db.Exec(
		`UPDATE Sites SET LastPing = $1
		  WHERE SiteId = $2`,
		p.TimeRequest,
		p.SiteID,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetFirstPing gets the earliest ping for a given site based on the recorded pings.
func (s *Site) GetFirstPing(db *sql.DB) (time.Time, error) {
	var firstPing *time.Time
	// The query is done with the aggregate MIN in the where clause because it doesn't
	//  work in the select clause with the timestamp type.
	err := db.QueryRow(`SELECT TimeRequest
		FROM Pings
		WHERE SiteID = $1 AND TimeRequest = (Select MIN(TimeRequest) FROM Pings
		WHERE SiteID = $1)`, s.SiteID).Scan(&firstPing)
	if err != nil {
		emptyTime := time.Time{}
		if err == sql.ErrNoRows {
			// This case isn't considered an error if it's an empty Ping table for this site.
			return emptyTime, nil
		}
		return emptyTime, err
	}

	return *firstPing, nil
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

// GetYTDReports gets reports for the active sites
func GetYTDReports(db *sql.DB) (map[string]Reports, error) {
	rows, err := db.Query(`
	SELECT Name, Month, SUM(AvgResponse) AS AvgResponse, SUM(PingsUp) As PingsUp, SUM(PingsDown) as PingsDown
	FROM(
	select Name, strftime("%m", timeRequest) as 'month', AVG(duration) as AvgResponse, count(*) as PingsUp, 0 as PingsDown
	     FROM pings INNER JOIN sites on sites.siteID = pings.siteID
		   WHERE httpstatuscode = 200 AND timeRequest > date('now', 'start of year')
		   group by strftime("%m", timeRequest), name
	UNION ALL
		   select Name, strftime("%m", timeRequest) as 'month', 0 as AvgResponse, 0 as PingsUp, count(*) as PingsDown
	       from pings INNER JOIN sites on sites.siteID = pings.siteID
		   WHERE httpstatuscode <> 200 AND timeRequest > date('now', 'start of year')
		   group by strftime("%m", timeRequest), name
	)
	group by name, month
	ORDER BY name, month`)
	if err != nil {
		return nil, err
	}

	var ytdReports map[string]Reports
	ytdReports = make(map[string]Reports)
	defer rows.Close()
	for rows.Next() {
		var Name string
		var Month int
		var AvgResponse float64
		var PingsUp int
		var PingsDown int
		err = rows.Scan(&Name, &Month, &AvgResponse, &PingsUp, &PingsDown)
		if err != nil {
			return nil, err
		}

		if _, ok := ytdReports[Name]; !ok {
			ytdReports[Name] = make([]Report, 12, 12)
		}
		ytdReports[Name][Month-1] = Report{AvgResponse: AvgResponse,
			PingsUp: PingsUp, PingsDown: PingsDown}
	}

	return ytdReports, nil
}
