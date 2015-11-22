package pinger_test

import (
	"strings"
	"testing"
	"time"

	"github.com/turnkey-commerce/go-ping-sites/notifier"
	"github.com/turnkey-commerce/go-ping-sites/pinger"
)

// TestNewPinger tests building the pinger object.
func TestNewPinger(t *testing.T) {
	pinger.CreatePingerLog()
	p := pinger.NewPinger(nil, pinger.GetSitesMock, pinger.RequestURLMock, notifier.SendEmailMock, notifier.SendSmsMock)

	if len(p.Sites) != 3 {
		t.Fatal("Incorrect number of sites returned in new pinger.")
	}

	results, err := pinger.GetLogContent()
	if err != nil {
		t.Fatal("Failed to get log results.", err)
	}

	if !strings.Contains(results, "SITE: Test, http://www.google.com") {
		t.Fatal("Failed to load first site.")
	}
	if !strings.Contains(results, "SITE: Test 2, http://www.github.com") {
		t.Fatal("Failed to load second site.")
	}
	if !strings.Contains(results, "SITE: Test 3, http://www.test.com") {
		t.Fatal("Failed to load third site.")
	}
}

// TestStartEmptySitesPinger starts up the pinger and then stops it after 1 second
func TestStartEmptySitesPinger(t *testing.T) {
	pinger.CreatePingerLog()
	p := pinger.NewPinger(nil, pinger.GetEmptySitesMock, pinger.RequestURLMock, notifier.SendEmailMock, notifier.SendSmsMock)
	p.Start()

	results, err := pinger.GetLogContent()
	if err != nil {
		t.Fatal("Failed to get log results.", err)
	}

	if !strings.Contains(results, "No active sites set up for pinging.") {
		t.Fatal("Failed to report empty sites.")
	}
}

// TestStartPinger starts up the pinger and then stops it after 10 seconds
func TestStartPinger(t *testing.T) {
	pinger.CreatePingerLog()
	p := pinger.NewPinger(nil, pinger.GetSitesMock, pinger.RequestURLMock, notifier.SendEmailMock, notifier.SendSmsMock)
	p.Start()
	time.Sleep(10 * time.Second)
	p.Stop()

	results, err := pinger.GetLogContent()
	if err != nil {
		t.Fatal("Failed to get log results.", err)
	}

	if !strings.Contains(results, "Client.Timeout") {
		t.Fatal("Failed to report timeout error.")
	}
	if !strings.Contains(results, "Test 3 Paused") {
		t.Fatal("Failed to report paused site.")
	}
	if !strings.Contains(results, "Error - HTTP Status Code") {
		t.Fatal("Failed to report bad HTTP Status Code.")
	}
	if !strings.Contains(results, "Will notify status change for Test 2 Site is now up.") {
		t.Fatal("Failed to report change in notification.")
	}
}
