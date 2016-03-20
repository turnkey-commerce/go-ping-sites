package pinger

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

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
func CreatePingerLog(logFilePath string, reinitializeLog bool) error {
	if logFilePath != "" {
		logFile = logFilePath
	}
	mustCreate := false
	if reinitializeLog {
		mustCreate = true
	} else {
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			mustCreate = true
		}
	}
	if mustCreate {
		_, err := os.Create(logFile)
		if err != nil {
			return err
		}
	}
	log.SetOutput(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    200, // MB
		MaxBackups: 3,
		MaxAge:     28, //days
	})
	return nil
}
