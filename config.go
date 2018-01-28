package main

// General configuration utilities for Glenda

import (
	"encoding/json"
	"os"
	"fmt"
)

// Stores config for current state
type Configuration struct {
	LastChange string
}

// Initializes current config (called once at start)
func (c *Configuration) Init() {
}

// Writes current config
func (c *Configuration) Write() (rerr error) {
	rerr = nil
	f, err := os.Open("./cfg/glenda.cfg")
	defer f.Close()
	if err != nil {
		fmt.Println("Error opening config, see: config.go")
		fmt.Printf("%s\n", err)
		rerr = err
	}

	e := json.NewEncoder(f)
	err = e.Encode(Config)
	if err != nil {
		fmt.Println("Error writing config, see: config.go")
		fmt.Printf("%s\n", err)
		rerr = err
	}
	return
}

// Reads current config into memory
func (c *Configuration) Read() {

}

// Set up config for the first time (if one doesn't exist)
func (c *Configuration) Setup() {
	// Test/make cfg folder


	// Test/make cfg


}

// Global variables are bad
var Config Configuration

