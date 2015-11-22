package pinger

import (
	"io/ioutil"
	"log"
	"os"
)

// LogFile is the name of the log file used to record the pinger activities.
var logFile = "pinger.log"

// GetLogContent reads the results of the log file for verification.
func GetLogContent() (string, error) {
	dat, err := ioutil.ReadFile(logFile)
	if err != nil {
		return "", err
	}
	results := string(dat)
	return results, nil
}

// CreatePingerLog creates the log file used by Pinger and Notifier
func CreatePingerLog() {
	pingerLog, err := os.Create(logFile)
	if err != nil {
		log.Fatal("Error creating pinger log", err)
	}
	log.SetOutput(pingerLog)
}
