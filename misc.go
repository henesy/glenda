package main

// Miscellaneous utilities for communicating with Mux

import (
//	"bitbucket.org/henesy/glenda/x/mux"
	"time"
)

// Service for communicating with Mux (starts once)
func CommMux() {
	for {
		select {
		default:
		}

		time.Sleep(500 * time.Millisecond)
	}
}

