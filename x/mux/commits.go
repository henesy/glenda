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
	for {
		for _, feed := range Config.Feeds {
			err := feed.Update()

			if err != nil {
				fmt.Println("Error in updating RSS feed, see: x/mux/commits.go")
				fmt.Printf("%s\n\n", err)
			}
			time.Sleep(1000 * time.Second)
		}
	}
}

/* Commands for Mux */

// Print last commit or 
func (m *Mux) Commits(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := "```\n"



	resp += "```\n"

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

// Subscribe to a feed
func (m *Mux) Subscribe(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	//http://code.9front.org/hg/plan9front/rss-log
	// TODO -- check if feed is already subscribed to
	resp := "```\n"
	
	feed, err := rss.Fetch("")

	if err != nil {
		fmt.Println("Error in reading RSS feed, see: x/mux/commits.go")
		fmt.Printf("%s\n\n", err)
	}
	
	// Might not be thread safe
	Config.Feeds = append(Config.Feeds, *feed)
	
	resp += "```\n"
	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}
