package main

// Miscellaneous utilities for communicating with Mux

import (
	"bitbucket.org/henesy/glenda/x/mux"
	"time"
)

// Service for communicating with Mux (starts once)
func CommMux() {
	for {
		select {
		case r := <-mux.GlendaChan:
			if r == "dump" {
				err := Config.Write()
				if err != nil {
					mux.MuxChan <- "dump failed"
				} else {
					mux.MuxChan <- "dump ok"
				}
			}
		default:
		}

		time.Sleep(500 * time.Millisecond)
	}
}

