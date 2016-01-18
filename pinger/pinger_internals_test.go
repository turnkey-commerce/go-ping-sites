package pinger

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type statusHandler int

func (h *statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(int(*h))
}
func TestIsInternetAccessibleFirstSiteBad(t *testing.T) {
	// Set up a non-running server and a good server
	statusOK := statusHandler(http.StatusOK)
	s1 := httptest.NewUnstartedServer(&statusOK)
	s2 := httptest.NewServer(&statusOK)
	defer s1.Close()
	defer s2.Close()

	result := isInternetAccessible(s1.URL, s2.URL)
	if !result {
		t.Error("Should pass on good second site.")
	}
}

func TestIsInternetAccessibleBothBad(t *testing.T) {
	// Set up both as non-running servers
	statusOK := statusHandler(http.StatusOK)
	s1 := httptest.NewUnstartedServer(&statusOK)
	s2 := httptest.NewUnstartedServer(&statusOK)
	defer s1.Close()
	defer s2.Close()

	result := isInternetAccessible(s1.URL, s2.URL)
	if result {
		t.Error("Should fail on both servers.")
	}
}

func TestIsInternetAccessibleSecondSiteBad(t *testing.T) {
	// Set up first OK and second as non-running server
	statusOK := statusHandler(http.StatusOK)
	s1 := httptest.NewServer(&statusOK)
	s2 := httptest.NewUnstartedServer(&statusOK)
	defer s1.Close()
	defer s2.Close()

	result := isInternetAccessible(s1.URL, s2.URL)
	if !result {
		t.Error("Should pass on first server.")
	}
}
