// Copyright 2015 Pajato Group Inc. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

// Package server handles the http traffic that provides the glue
// between CLI operations and deamon processing.
package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// State encapsulates relevant server state.
type State struct {
	// The version established when the server was setup.
	Version string
	// The Process ID for the process running the server.
	Pid int
}

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var (
	// The set of operations supported by this server.
	routes = []route{
		{"AuthenticationCallback", "GET", "/authenticationCallback", authenticationHandler},
		{"PreferencesPage", "GET", "/preferences", preferencesHandler},
		{"StatusPage", "GET", "/status", statusHandler},
		{"PingTest", "GET", "/ping", pingHandler},
		{"ShowLogPage", "GET", "/showlog", showLogHandler},
	}

	state State
)

// SetupServer sets up the local server that handles preferences and
// authentication, among other things.
func SetupServer(port string, version string) {
	// Capture the arg (server) and create the router using the route
	// table.
	log.Printf("Setting up server to run at url: %s\n", port)
	state = State{version, os.Getpid()}
	r := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		log.Printf("Setting up route for %s\n", route.Name)
		h := getLoggingHandler(route.HandlerFunc, route.Name)
		r.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(h)
	}
	//router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets/"))))
	log.Print("Starting server...\n")
	log.Fatal(http.ListenAndServe(port, r))
}

// Obtains the authentication key information provided by Google
// giving access to Google Drive for the logged in user.
func authenticationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "tbd")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, state.Version+":"+strconv.Itoa(state.Pid))
}

func preferencesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "tbd")
}

func showLogHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "tbd")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "tbd")
}

// The getLoggingHandler() function ensures that each request is logged and
// returns a handler that executes a given, named handler.
func getLoggingHandler(handler http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		log.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))
	})
}
