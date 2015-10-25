// Copyright 2015 Pajato Group Inc. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

// GDsync Documentation
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/pajato/gdsync/server"
)

type handler func()

const (
	// ServerURL defines the port for the server handling web traffic.
	serverPort = ":8089"

	// The Program version
	programVersion = "0.1.0"
)

var (
	// Path points to a string specifying the full OS path to the file
	// containing stored preferences.
	path = *flag.String("c", os.Getenv("HOME")+"/.gdsync",
		"The OS path to the file containing GDsync preferences.")

	// Subcommand map associating a name with a handler.
	m = map[string]handler{
		"adduser":     addUser,
		"showlog":     showLog,
		"startdaemon": startDaemonMaybe,
		"startserver": startServerMaybe,
	}

	// A init() controlled variable (defaults to false) that controls
	// which ListenAndServer variant to use.  Provided to enhance test
	// coverage reporting.
	test bool
)

// Determine if the daemon is running.  If not, it will be started.
// If the daemon is running, the arguments will be processed and the
// process will terminate.
func main() {
	// Grab the flags to establish the preferences file and determine
	// if the daemon is running in another process.
	flag.Parse()
	processSubCommand(flag.Arg(0))
}

// Process the args outside of main so that testing code can simulate
// the passing of arguments to the command line.
func processSubCommand(subcommand string) {
	// ...
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.Printf("Preferences file path is '%s'.\n", path)
	c := strings.ToLower(subcommand)
	f := m[c]
	if f != nil {
		f()
	} else {
		log.Printf("Invalid command: %v.  Ignoring.\n", c)
	}
}

func addUser() {
	// tbd
	ensureServerIsRunning()
}

func ensureServerIsRunning() {
	startDaemonMaybe()
}

func showLog() {
	// tbd
}

func killDaemon(pid int) {
	s := "Could not kill the daemon process.  Failed with: %v.  Aborting.\n"
	if p, err := os.FindProcess(pid); err != nil {
		log.Fatalf(s, err)
	} else {
		if err := p.Kill(); err != nil {
			log.Fatalf(s, err)
		}
	}
}

func startDaemonMaybe() {
	if version, pid, err := ping(); err != nil {
		log.Printf("Could not ping the server.  Failed with: %v.\n", err)
		log.Print("Starting up daemon...")
		startDaemonAndWait()
	} else {
		log.Print("Daemon is running.  Checking server state.")
		if version != programVersion {
			log.Printf("Server has wrong version: %s.  Restarting daemon.\n", version)
			killDaemon(pid)
			startDaemonAndWait()
		}
	}
}

func startServerMaybe() {
	// This is a daemon invocation command.
	if v, p, err := ping(); err != nil {
		log.Printf("Could not ping the server.  Failed with: %v\n.", err)
		log.Print("Starting up server...")
		server.SetupServer(serverPort, programVersion)
	} else {
		log.Printf("Server is already running with version {%s} on port {%d}.\n", v, p)
	}
}

func ping() (version string, pid int, err error) {
	// Ping the server running in the daemon returning the version,
	// process id and error information discerned when processing the
	// request.
	var r *http.Response
	if r, err = http.Get("http://localhost" + serverPort + "/ping"); err == nil {
		defer r.Body.Close()
		var body []byte
		if body, err = ioutil.ReadAll(r.Body); err == nil {
			v := string(body[:len(body)])
			log.Printf("Ping value from server: {%s}.\n", v)
			args := strings.Split(v, ":")
			if n := len(args); n != 2 {
				log.Printf("Wrong number of parameters returned by the server: %d.\n", n)
				return
			}
			version = args[0]
			if pid, err = strconv.Atoi(args[1]); err != nil {
				log.Printf("Cannot convert PID {%s} to an int.\n", args[1])
				return
			}
			if version != programVersion {
				log.Printf("The server is running the wrong version: %v.\n", version)
				return
			}
			log.Println("Ping completed successfully.")
			return
		}
	}
	log.Printf("The server is not running.  The failure is: %v", err)
	return
}

func startDaemonAndWait() {
	cmd := exec.Command("gdsync", "startserver", programVersion)
	if err := cmd.Start(); err != nil {
		log.Fatalf("Could not start the daemon.  Failed with: %v.\n", err)
	} else {
		// Wait a few seconds for the daemon to start the server.
		var c int
		for c = 30; c > 0; c-- {
			time.Sleep(100 * time.Millisecond)
			if _, _, err := ping(); err == nil {
				break
			}
		}
		if c == 0 {
			// Give up and abort.
			log.Fatal("Could not start the daemon.  Aborting.")
		} else {
			// The daemon is started.  Continue.
			log.Print("Google Drive sync daemon started.  Continuing.")
		}
	}
}
