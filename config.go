package main

// General configuration utilities for Glenda

import (
	"encoding/json"
	"os"
	"fmt"
	"strings"
)

// Stores config for current state
type Configuration struct {
	LastChange string
}

// Initializes current config (called once at start) ­ just .Read()?
func (c *Configuration) Init() {
	c.Read()
}

// Writes current config
func (c *Configuration) Write() (rerr error) {
	WRT: 
	rerr = nil
	f, err := os.OpenFile("./cfg/glenda.cfg", os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			Config.Setup()
			goto WRT
			return
		} else {
			fmt.Println("Error opening config (w), see: config.go")
			fmt.Printf("%s\n", err)
			rerr = err
		}
	} else {
		e := json.NewEncoder(f)
		err = e.Encode(Config)
		if err != nil {
			fmt.Println("Error writing config, see: config.go")
			fmt.Printf("%s\n", err)
			rerr = err
		}
	}
	
	return
}

// Reads current config into memory
func (c *Configuration) Read() (rerr error) {
	f, err := os.Open("./cfg/glenda.cfg")
	if err != nil {
		fmt.Println("Error opening config (r), see: config.go")
		fmt.Printf("%s\n", err)
		rerr = err
	}
	
	d := json.NewDecoder(f)
	err = d.Decode(&Config)
	if err != nil {
		fmt.Println("Error reading config, see: config.go")
		fmt.Printf("%s\n", err)
		rerr = err
	}
	
	return
}

// Set up config for the first time (if one doesn't exist)
func (c *Configuration) Setup() {
	err := os.Mkdir("./cfg", 0774)
	if err != nil {
		fmt.Println("Error in making cfg dir, see: config.go")
		fmt.Println(err)
	}
	
	_, err = os.Create("cfg/glenda.cfg")
	if err != nil {
		fmt.Println("Error in making cfg file, see: config.go")
		fmt.Println(err)
	}
}

// Global variables are bad
var Config Configuration

