package pinger

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

// TestRequestURL tests the production implementation of the RequestURL code by
// requesting an actual site.
func TestRequestURL(t *testing.T) {
	content, responseCode, responseTime, err := RequestURL("http://www.example.com", 60)
	if err != nil {
		t.Error("Request URL retrieval error", err)
	}

	if !strings.Contains(content, "<title>Example Domain</title>") {
		t.Error("Request URL response code error", responseCode)
	}

	if responseCode != 200 {
		t.Error("Request URL response code error", responseCode)
	}
}

// TestRequestURLError tests the error handling of the production implementation
// of the RequestURL code by requesting a bogus site that will throw an error.
func TestRequestURLError(t *testing.T) {
	_, _, _, err := RequestURL("http://www.examplefoobar.com", 5)
	if err == nil {
		t.Error("Bad URL should throw error")
	}
}

// TestRequestInternetAccessError tests the error handling of the production
// implementation by simulating Internet access error with a bad test sites.
func TestRequestInternetAccessError(t *testing.T) {
	site1 = "http://www.examplefoobar.com"
	site2 = "http://www.examplefoobar2.com"
	_, _, _, err := RequestURL("http://www.examplefoobar.com", 5)
	if err == nil {
		t.Error("Bad URL and test sites should throw error")
	}
	if _, ok := err.(InternetAccessError); !ok {
		t.Error("Bad URL and test sites should identify as Internet access error.")
	}
}
