package mux

import (
	"github.com/bwmarrin/discordgo"
	"github.com/SlyMarbo/rss"
	"fmt"
	"time"
)

/* internal functions */

// Listen on the rss feed
func Listener() {
	feed, err := rss.Fetch("http://code.9front.org/hg/plan9front/rss-log")
	if err != nil {
		fmt.Println("Error in reading RSS feed, see: x/mux/commits.go")
		fmt.Printf("%s\n\n", err)
	}

	time.Sleep(1000 * time.Second)

	err = feed.Update()
	if err != nil {
		fmt.Println("Error in updating RSS feed, see: x/mux/commits.go")
		fmt.Printf("%s\n\n", err)
	}
}

// Load JSON Configuration File (move me?)


/* Commands for Mux */

// Print last commit or 
func (m *Mux) Commits(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := "```\n"



	resp += "```\n"

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}
