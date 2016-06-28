package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

var configFile = "config.toml"

// Settings contains the settings for SMTP, Twilio, and Website from the config.toml file.
var Settings struct {
	SMTP struct {
		EmailAddress string `valid:"-"`
		Password     string `valid:"-"`
		Server       string `valid:"-"`
		Port         string `valid:"int,required"`
	}
	Twilio struct {
		AccountSid string `valid:"-"`
		AuthToken  string `valid:"-"`
		Number     string `valid:"-"`
	}
	Website struct {
		HTTPPort    string `valid:"int,required"`
		CookieKey   string `valid:"ascii,required"`
		SecureHTTPS bool   `valid:"bool"`
	}
}

func init() {
	if _, err := toml.DecodeFile(configFile, &Settings); err != nil {
		log.Fatal(err)
	}
}
