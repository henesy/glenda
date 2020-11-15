package mux

// General configuration utilities for Glenda

import (
	"encoding/json"
	"os"
	"fmt"
	"strings"
	"github.com/SlyMarbo/rss"
	"github.com/bwmarrin/discordgo"
	"time"
)


// Global variables are bad
var Config Configuration
var Session *discordgo.Session

// Stores subscription information for channels related to RSS feeds
type Subscription struct {
	ChanID	string
	SubID	int
}

// Stores a Feed and a list of the last n commits; n=3
type Feed struct {
	Feed		rss.Feed
	Recent	[]string
}

type Feeder struct {
	Feeds	[]Feed
	Subs		[]Subscription
}

// Stores config for current state
type Configuration struct {
	Db			map[string]string
	Feeder
	Reminders []Reminder
}

// Initializes current config (called once at start) ­ just .Read()?
func (c *Configuration) Init(s *discordgo.Session) {
	c.Read()
	Session = s
	for id, _ := range Config.Feeds {
		// maybe only do at init step?
		str := Config.Feeds[id].Feed.UpdateURL
		feed, _ := rss.Fetch(str)
		if feed != nil {
			Config.Feeds[id].Feed = *feed
		} else {
			fmt.Println("Failed to fetch feed: ", id)
		}
	}
	
	if c.Db == nil {
		c.Db = map[string]string {
		"name":	"glenda.cfg",
		"dir": "./cfg",
		}
	}
	
	go Listener()
}

// Writes current config
func (c *Configuration) Write() (rerr error) {
	WRITE:
	rerr = nil
	f, err := os.OpenFile(c.Db["dir"] + "/" + c.Db["name"], os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			// Create and try again
			// danger: this can go infinite
			Config.Setup()
			time.Sleep(5 * time.Millisecond)
			goto WRITE
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
	READ:
	f, err := os.Open(c.Db["dir"] + "/" + c.Db["name"])
	defer f.Close()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			// danger: this can go infinite
			Config.Setup()
			goto READ
		} else {
			fmt.Println("Error opening config (r), see: config.go")
			fmt.Printf("%s\n", err)
			rerr = err
		}
	} else {
		d := json.NewDecoder(f)
		err = d.Decode(&Config)
		if err != nil {
			fmt.Println("Error reading config, see: config.go")
			fmt.Printf("%s\n", err)
			rerr = err
			Config.Write()
		}
	}

	return
}

// Set up config for the first time (if one doesn't exist)
func (c *Configuration) Setup() {
	err := os.Mkdir(c.Db["dir"], 0774)
	if err != nil {
		fmt.Println("Error in making cfg dir, see: config.go")
		fmt.Println(err)
	}
	
	_, err = os.Create(c.Db["dir"] + "/" + c.Db["name"])
	if err != nil {
		fmt.Println("Error in making cfg file, see: config.go")
		fmt.Println(err)
	}
}
