// testing...
package pref

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/pajato/bus"
)

// Test that when no preferences are found at a given location,
// ensured by the test, a set of zero-valued preferences are then
// created successfully.
func TestNoPreferencesSuccess(t *testing.T) {
	log.Println("Testing preferences expecting success.")
	// Create a bus object, ensure that the preferences file does not
	// exist, create the default preferences and save it. Finally,
	// verify the default values read in from the preferences file.
	path := "/tmp/pref-test/.prefs"
	conditionallyRemove(t, path)
	createAndSaveDefaults(t, path)
	prefs := loadDefaults(t, path)
	verify(t, prefs, "", os.Getenv("HOME")+"/GoogleDrive", false, 0)
}

// Test that when no preferences are found at a given location,
// ensured by the test, a set of zero-valued preferences cannot then
// be created.
func TestNoPreferencesFailWithMissingDirectory(t *testing.T) {
	log.Println("Testing preferences expecting failure due to missing directory.")
	// Create a bus object, ensure that the preferences file does not
	// exist, create the default preferences and save it. Finally,
	// verify the default values read in from the preferences file.
	path := "/there_is_no_such_directory/pref-test/.prefs"
	failToCreatePrefs(t, path)
}

// Test that when a valid but inacessible directory is specified, no
// preferences can be created.
func TestNoPreferencesFailWithInaccessibleDirectory(t *testing.T) {
	log.Println("Testing preferences expecting failure due to inacessible directory.")
	// Create a bus object, ensure that the preferences file does not
	// exist, create the default preferences and save it. Finally,
	// verify the default values read in from the preferences file.
	path := "/var/backups/.prefs"
	failToCreatePrefs(t, path)
}

// Test that when no preferences are found at a given location,
// ensured by the test, a set of zero-valued preferences cannot then
// be created.
func TestNoPreferencesFailWithBadFile(t *testing.T) {
	log.Println("Testing preferences expecting failure due to bad path.")
	// Create a bus object, ensure that the preferences file does not
	// exist, create the default preferences and save it. Finally,
	// verify the default values read in from the preferences file.
	path := "/.prefs"
	failToCreatePrefs(t, path)
}

// Test that the server handles an authenticaion key...
func TestAuthenticationKey(t *testing.T) {
	log.Println("Testing the authentication key response.")
	path := "./test-data/.prefs"
	bus := bus.New()
	if _, err := Launch(path, bus); err != nil {
		t.Fatalf("Failed to setup authentication key test.  Failed with {%v}.\n", err)
	}
	if resp, err := http.Get(ServerURL + "/authenticationCallback"); err != nil {
		// Deal with errors
		t.Fatalf("Could not get the authentication key!  Failed with %v.\n", err)
	} else {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Fatalf("Could not read the response body!  Failed with %v.\n", err)
		} else {
			n := len(body)
			s := string(body[:n])
			log.Printf("The authentication key is {%v}.\n", s)
		}
	}
}

func conditionallyRemove(t *testing.T, path string) {
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
		_, err := os.Open(path)
		if err == nil {
			t.Error("The preferences file cannot be removed.  Test cannot be performed.\n")
		}
	}
}

func failToCreatePrefs(t *testing.T, path string) {
	bus := bus.New()
	if _, err := Launch(path, bus); err == nil {
		t.Errorf("The prefs file was created when it should not have been!")
	}
}

func createAndSaveDefaults(t *testing.T, path string) {
	bus := bus.New()
	if _, err := Launch(path, bus); err == nil {
		if _, err := os.Open(path); err != nil {
			t.Errorf("The prefs file was not created, failed with: %s\n", err)
		}
	} else {
		t.Errorf("The prefs file was not created, do to a launch error: {%s}.\n", err)
	}
}

func loadDefaults(t *testing.T, path string) Pref {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		t.Errorf("The prefs file cannot be opened: {%s}\n", err)
	}
	decoder := gob.NewDecoder(file)
	prefs := Pref{}
	err = decoder.Decode(&prefs)
	format := "User is: {%s}; Path is: {%s}; SyncAll is: {%s}; Number of folders to sync is: {%d}.\n"
	log.Printf(format, prefs.User, prefs.Base, prefs.SyncAll, len(prefs.SyncFolders))
	return prefs
}

func verify(t *testing.T, prefs Pref, user string, base string, syncAll bool, size int) {
	message := ""
	format := "The %s preference value is: %s.  %s was expected.\n"
	if prefs.User != user {
		message += fmt.Sprintf(format, "user", prefs.User, user)
	}
	if prefs.Base != base {
		message += fmt.Sprintf(format, "base", prefs.Base, base)
	}
	if prefs.SyncAll {
		message += fmt.Sprintf(format, "syncAll", prefs.SyncAll, syncAll)
	}
	if len := len(prefs.SyncFolders); len != size {
		format = "The syncFolders size is: %d.  %d was expected.\n"
		message += fmt.Sprintf(format, len, size)
	}
	if len(message) != 0 {
		t.Errorf("The preferences are not as expected.\n%s", message)
	}
}
