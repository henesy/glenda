package mux

// General configuration utilities for Glenda

import (
	"encoding/json"
	"os"
	"fmt"
	"strings"
	"github.com/SlyMarbo/rss"
	"github.com/bwmarrin/discordgo"
	"container/list"
)


// Global variables are bad
var Config Configuration
var Feeds	Feeder
var Rems *list.List
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

}

// Initializes current config (called once at start) ­ just .Read()?
func (c *Configuration) Init(s *discordgo.Session) {
	Rems = list.New()
	c.Read()
	Session = s
	for id, _ := range Feeds.Feeds {
		// maybe only do at init step?
		str := Feeds.Feeds[id].Feed.UpdateURL
		feed, _ := rss.Fetch(str)
		if feed != nil {
			Feeds.Feeds[id].Feed = *feed
		} else {
			fmt.Println("Failed to fetch feed: ", id)
		}
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

// Writes the Feed config file
func (c *Feeder) Write() (rerr error) {
	FWRT:
	rerr = nil
	f, err := os.OpenFile("./cfg/feeds.cfg", os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			// danger: this can go infinite
			Config.Setup()
			goto FWRT
		} else {
			fmt.Println("Error opening feeds config (w), see: config.go")
			fmt.Printf("%s\n", err)
			rerr = err
		}
	} else {
		e := json.NewEncoder(f)
		err = e.Encode(Feeds)
		if err != nil {
			fmt.Println("Error writing feeds config, see: config.go")
			fmt.Printf("%s\n", err)
			rerr = err
		}
	}

	return
}

// Writes the Reminders config file
func RemindersWrite() (rerr error) {
	RWRT:
	rerr = nil
	f, err := os.OpenFile("./cfg/reminders.cfg", os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			// danger: this can go infinite
			Config.Setup()
			goto RWRT
		} else {
			fmt.Println("Error opening reminders config (w), see: config.go")
			fmt.Printf("%s\n", err)
			rerr = err
		}
	} else {
		e := json.NewEncoder(f)
		err = e.Encode(*Rems)
		if err != nil {
			fmt.Println("Error writing reminders config, see: config.go")
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
	f.Close()
	
	// Feeds
	f, err = os.Open("./cfg/feeds.cfg")
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			// danger: this can go infinite
			Config.Setup()
			goto RD
		} else {
			fmt.Println("Error opening feeds config (r), see: config.go")
			fmt.Printf("%s\n", err)
			rerr = err
		}
	} else {
		d := json.NewDecoder(f)
		err = d.Decode(&Feeds)
		if err != nil {
			fmt.Println("Error reading feeds config, see: config.go")
			fmt.Printf("%s\n", err)
			rerr = err
			Feeds.Write()
		}
	}
	f.Close()
	
	// Reminders
	f, err = os.Open("./cfg/reminders.cfg")
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			// danger: this can go infinite
			Config.Setup()
			goto RD
		} else {
			fmt.Println("Error opening reminders config (r), see: config.go")
			fmt.Printf("%s\n", err)
			rerr = err
		}
	} else {
		d := json.NewDecoder(f)
		err = d.Decode(Rems)
		if err != nil {
			fmt.Println("Error reading reminders config, see: config.go")
			fmt.Printf("%s\n", err)
			rerr = err
			RemindersWrite()
		}
	}
	f.Close()

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
	
	// Feeds
	_, err = os.Create("cfg/feeds.cfg")
	if err != nil {
		fmt.Println("Error in making feed cfg file, see: config.go")
		fmt.Println(err)
	}
	
	// Reminders
	_, err = os.Create("cfg/reminders.cfg")
	if err != nil {
		fmt.Println("Error in making reminders cfg file, see: config.go")
		fmt.Println(err)
	}
}
