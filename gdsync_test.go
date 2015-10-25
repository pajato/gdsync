// Copyright 2015 Pajato Group Inc. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

// Package bus is a subscription based (Publisher/Subscriber pattern)
// communication subsystem based loosely on the Java/Android Otto
// library, which in turn is based on Guava.
package main

import (
	"log"
	"testing"
)

func TestEmptySubCommand(t *testing.T) {
	log.Println("Testing 'gdsync'")
	processSubCommand("")
}

func TestAddUserSubCommand(t *testing.T) {
	log.Println("Testing 'gdsync addUser'")
	processSubCommand("addUser")
}
