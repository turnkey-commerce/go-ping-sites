package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

var configFile = "config.toml"

// Settings stores the configuration for the SMTP email and Twilio SMS accounts.
var Settings struct {
	SMTP struct {
		EmailAddress string
		Password     string
		Server       string
		Port         string
	}
	Twilio struct {
		AccountSid string
		AuthToken  string
		Number     string
	}
	Website struct {
		HTTPPort  string
		CookieKey string
	}
}

func init() {
	if _, err := toml.DecodeFile(configFile, &Settings); err != nil {
		log.Fatal(err)
	}
}
