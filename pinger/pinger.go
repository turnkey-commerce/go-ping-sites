package pinger

import (
	"net/http"

	"github.com/turnkey-commerce/go-ping-sites/database"
)

// Pinger does the HTTP pinging of the sites that are retrieved from the DB.
type Pinger struct {
	sites database.Sites
	http  *http.Client
}
