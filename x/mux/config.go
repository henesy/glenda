package mux

// General configuration utilities for Glenda

import (
	"encoding/json"
	"os"
	"fmt"
	"strings"
	"github.com/SlyMarbo/rss"
	"github.com/bwmarrin/discordgo"
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

// Stores config for current state
type Configuration struct {
	Feeds	[]Feed
	Subs		[]Subscription
}

// Initializes current config (called once at start) ­ just .Read()?
func (c *Configuration) Init(s *discordgo.Session) {
	c.Read()
	Session = s
	for id, _ := range Config.Feeds {
		// maybe only do at init step?
		str := Config.Feeds[id].Feed.UpdateURL
		feed, _ := rss.Fetch(str)
		Config.Feeds[id].Feed = *feed
	}
	go Listener()
}

// Writes current config
func (c *Configuration) Write() (rerr error) {
	WRT:
	rerr = nil
	f, err := os.OpenFile("./cfg/glenda.cfg", os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			// danger: this can go infinite
			Config.Setup()
			goto WRT
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
	RD:
	f, err := os.Open("./cfg/glenda.cfg")
	defer f.Close()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			// danger: this can go infinite
			Config.Setup()
			goto RD
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
