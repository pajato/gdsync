// Copyright 2015 Pajato Group Inc. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

// Package pref package handles gdsync preferences and provides a GUI for
// the app.
package pref

import (
	"encoding/gob"
	"log"
	"os"

	pp "path"

	"github.com/pajato/bus"
	"github.com/pkg/browser"
)

// A Pref object provides a container for User modifiable variables.
type Pref struct {
	User        string
	Base        string
	SyncAll     bool
	SyncFolders []string
}

// Launch starts the GDsync preference task, using a file based
// preferences store rooted at the given relative or absolute path.
// The directory and/or file components of path are created as
// necessary.  If there are file or directory creation or access
// errors, err will return a *PathError type value and prefs will be
// nil, otherwise prefs will contain a valid (possibly default) Pref
// type value and err will be nil.
func Launch(path string, b bus.Bus) (prefs *Pref, err error) {
	// Deal with the preferences file: create some zero valued
	// preferences if they do not already exist.  Then unmarshall the
	// preferences file.
	if _, err = os.Stat(path); err != nil {
		// There is an issue with the given path.  Deal with it.
		if os.IsNotExist(err) {
			// The preferences file does not exist.  Create it now and
			// seed it with defaults.
			if err = createAndMarshallDefaults(path); err != nil {
				// The default preferences could not be created.  Return the error.
				format := "Failed to create/marshall the preferences file.  The error is: %v.\n"
				log.Printf(format, err)
				return
			}
		} else {
			// There is some unexpected and unanticipated path problem so abort.
			return
		}
	}

	// Open and read in the preferences.
	if prefs, err = openAndUnmarshallPrefs(path); err == nil {
		// Notify subscribers that the preferences file is ready.
		bus.Payload = ??
	}
	return
}

// Create default preferences and marshall them to a preferences store
// at the given path.
func createAndMarshallDefaults(path string) (err error) {
	// Split the path into its base file and directory parts.  Then
	// ensure that the directory exists.
	if dir, _ := pp.Split(path); len(dir) > 0 {
		// Validate the directory.
		if _, err = os.Stat(dir); err != nil {
			// A problem exists either because the directory needs to
			// be created or is inaccessible.
			if os.IsNotExist(err) {
				// The directory does not exist.  Create it now.
				if err = os.MkdirAll(dir, 0755); err != nil {
					// Failed due to inability to create the directory.
					return
				}
			}
		}
	}

	// Now create the preferences store.
	var file *os.File
	if file, err = os.Create(path); err != nil {
		// Aborting since the preferences file cannot be created.
		return
	}
	defer file.Close()
	log.Printf("Preferences file: %s created.\n", file.Name())
	encoder := gob.NewEncoder(file)
	prefs := Pref{Base: os.Getenv("HOME") + "/GoogleDrive"}
	err = encoder.Encode(prefs)
	return
}

func openAndUnmarshallPrefs(path string) (prefs *Pref, err error) {
	var file *os.File
	if file, err = os.Open(path); err != nil {
		// Abort since the prefs file cannot be opened.
		return
	}
	defer file.Close()
	log.Printf("Using file %s for preferences.\n", file.Name())
	decoder := gob.NewDecoder(file)
	prefs = &Pref{}
	err = decoder.Decode(prefs)
	format := "User is: {%s}; Path is: {%s}; SyncAll is: {%s}; Number of folders to sync is: {%d}.\n"
	log.Printf(format, prefs.User, prefs.Base, prefs.SyncAll, len(prefs.SyncFolders))
	return
}

// ProcessPrefs ensures that the preferences have reasonable values.
func (prefs Pref) ProcessPrefs() {
	// Start processing the "user" and other preference.
	if len(prefs.User) == 0 {
		log.Println("Need to establish the User ... start the web interfaces.")
		if err := browser.OpenURL("http://golang.org"); err != nil {
			log.Printf("Could not open browser.  Failed with: %s\n", err)
		}
	}
	log.Printf("Syncing on behalf of User: %s\n", prefs.User)
}
