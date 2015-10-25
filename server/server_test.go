// testing...
package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

const (
	// ServerURL defines the protocol, IP address and port for the
	// server handling web traffic.
	serverPort = ":8089"
)

// TestMain sets up the server under test to run in a goroutine.
func TestMain(m *testing.M) {
	go SetupServer(serverPort, "0.0.0")
	os.Exit(m.Run())
}

// Test that the server handles the authenticaion key...
func TestAuthentication(t *testing.T) {
	testPage(t, "/authenticationCallback")
}

// Test that the server handles the landing page...
func TestPreferencesPage(t *testing.T) {
	testPage(t, "/preferences")
}

// Test that the server handles an authenticaion key...
func TestPing(t *testing.T) {
	testPage(t, "/ping")
}

// Test that the server handles an authenticaion key...
func TestShowLogPage(t *testing.T) {
	testPage(t, "/showlog")
}

// Test that the server handles an authenticaion key...
func TestStatusPage(t *testing.T) {
	testPage(t, "/status")
}

func testPage(t *testing.T, page string) {
	log.Printf("Testing the %s response.", page)
	//SetupServer(serverPort, "0.0.0", true)
	if resp, err := http.Get("http://localhost" + serverPort + page); err != nil {
		// Deal with errors
		t.Fatalf("Could not get the response!  Failed with %v.\n", err)
	} else {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Fatalf("Could not read the response body!  Failed with %v.\n", err)
		} else {
			n := len(body)
			s := string(body[:n])
			log.Printf("The response is {%v}.\n", s)
		}
	}
}
