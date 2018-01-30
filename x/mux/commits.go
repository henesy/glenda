package mux

import (
	"github.com/bwmarrin/discordgo"
	"github.com/SlyMarbo/rss"
	"fmt"
	"time"
	"strings"
	"strconv"
)

/* internal functions */

// Listen on the rss feed
func Listener() {
	for {
		for _, feed := range Config.Feeds {
			fmt.Print("Printing feed: ", feed, "\n")
			err := feed.Update()

			if err != nil {
				fmt.Println("Error in updating RSS feed, see: x/mux/commits.go")
				fmt.Printf("%s\n\n", err)
			}
			time.Sleep(10 * time.Minute)
		}
		time.Sleep(10 * time.Minute)
	}
}

/* Commands for Mux */

// Print last commit for a given feed (by subscribed id int [See: List(...)])
func (m *Mux) Last(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := ".\n"

	id, _ := strconv.Atoi(ctx.Fields[len(ctx.Fields) -1])
	if id >= 0 && id < len(Config.Feeds) {
		resp += "**" + Config.Feeds[id].Title + ": **" + "\n"
		resp += Config.Feeds[id].Items[ len(Config.Feeds[id].Items) -1].Date.String() + "\n\n"
		resp += "`" + Config.Feeds[id].Items[ len(Config.Feeds[id].Items) -1].Title + "`" + "\n"
		//resp += Config.Feeds[id].Items[ len(Config.Feeds[id].Items) -1].Summary + "\n"
		//resp += Config.Feeds[id].Items[ len(Config.Feeds[id].Items) -1].Content + "\n"
		resp += "\n" + Config.Feeds[id].Items[ len(Config.Feeds[id].Items) -1].Link + "\n"
	} else {
		resp += "Denied fetch. Invalid stream id, see: list command"
	}

	resp += "\n"

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

// List all subscribed feeds
// Print last commit for a given feed
func (m *Mux) List(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := "```\n"

	for p, v := range Config.Feeds {
		resp += strconv.Itoa(p) + ": " + v.Title + "\n"
	}

	resp += "```\n"

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

// Subscribe to a feed
func (m *Mux) Subscribe(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	//http://code.9front.org/hg/plan9front/rss-log
	resp := "```\n"
	// URL to feed should be last item
	url := ctx.Fields[len(ctx.Fields) -1]
	fmt.Println("Proposed subscribe for: ", url)
	
	for _, v := range Config.Feeds {
		if strings.Contains(url, v.Link) {
			resp += "Denied! Feed already subscribed to."
			resp += "```\n"
			ds.ChannelMessageSend(dm.ChannelID, resp)
			return
		}
	}
	
	feed, err := rss.Fetch(url)

	if err != nil {
		fmt.Println("Error in reading RSS feed, see: x/mux/commits.go")
		fmt.Printf("%s\n\n", err)
		resp += "Denied! Could not parse feed."
		resp += "```\n"
		ds.ChannelMessageSend(dm.ChannelID, resp)
		return
	}
	
	// Might not be thread safe
	Config.Feeds = append(Config.Feeds, *feed)
	resp += "Subscribed."
	
	resp += "```\n"
	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}
