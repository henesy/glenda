package mux

// General configuration utilities for Glenda

import (
	"encoding/json"
	"os"
	"fmt"
	"strings"
	"errors"
	"strconv"
	"time"
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
	c.Db = map[string]string {
		"owner":		"188698402727526400",			// Henesy
		"name":		"glenda.cfg",					// Config file name
		"dir":			"./cfg",						// Config dir
		"delay":		"20",							// Minutes
		"mansite":		"http://man.postnix.pw/9front/",	// Henesy's
		
		"lookman":	"./x/mux/lookman/lookman",		// Vendored
		"sig":		"./x/mux/man/sig",			// ^
		"gendex":		"./gendex",					// ^
		"bullshit":		"./x/mux/misc/bullshit",			// ^
		
		"extrafortunes":	"../sys/games/lib/fortunes",			// Unix
		"fortunes":		"/usr/share/mirror/plan9front/lib/",	// Plan 9
		
		}
		
	// Regularly dump the config
	go func() {
		delay, err := strconv.Atoi(c.Db["delay"])
		if err != nil {
			fmt.Println("bad delay value, using default")
			delay = 20
		}
		
		for {
			select {
			case <-dumpChan:
				dump()

			case <- time.After(time.Duration(delay) * time.Minute):
				dump()
			}
			time.Sleep(5 * time.Millisecond)
		}
	}()

	err := c.Read()
	if err != nil {
		fmt.Println("read cfg failed: -", err)
		
		c.Setup()
	}
	
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
	
	go Listener()
}

// Writes current config to file
func (c *Configuration) Write() error {
	path := c.Db["dir"] + "/" + c.Db["name"]
	
	if path == "/" {
		return errors.New("'dir' and 'name' must be in config")
	}
	
	var f *os.File
	var err error
	
	WRITE:
	f, err = os.OpenFile(path, os.O_RDWR, 0666)
	defer f.Close()
	
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			// Create files
			err = Config.Setup()
			if err != nil {
				// We have a creation problem
				return err
			}
			
			// Try again
			goto WRITE
			
		} else {
			fmt.Println("Error opening config (w), see: config.go -", err)
			return err
		}
	} else {
		// Serialize
		e := json.NewEncoder(f)
		err = e.Encode(Config)
		if err != nil {
			fmt.Println("Error writing config (w), see: config.go -", err)
			defer Config.Write()
			
			return err
		}
	}

	return nil
}

// Reads current config into memory
func (c *Configuration) Read() error {
	path := c.Db["dir"] + "/" + c.Db["name"]
	
	if path == "/" {
		return errors.New("'dir' and 'name' must be in config")
	}

	var f *os.File
	var err error
	
	f, err = os.Open(path)
	defer f.Close()
	
	if err != nil {
		return err
	}
	
	// De-serialize
	d := json.NewDecoder(f)
	err = d.Decode(&Config)
	if err != nil {
		fmt.Println("Error reading config (r), see: config.go -", err)
		
		// Overwrite since the file is bad
		defer Config.Write()
		
		return err
	}

	return nil
}

// Set up config for the first time (if one doesn't exist)
func (c *Configuration) Setup() error {
	err := os.Mkdir(c.Db["dir"], 0774)
	if err != nil {
		fmt.Println("Error in making cfg dir, see: config.go -", err)
		
		if !strings.Contains(err.Error(), "exists") {
			return err
		}
	}
	
	_, err = os.Create(c.Db["dir"] + "/" + c.Db["name"])
	if err != nil {
		fmt.Println("Error in making cfg file, see: config.go -", err)
		
		if !strings.Contains(err.Error(), "exists") {
			return err
		}
	}
	
	return nil
}
