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
		for p, feed := range Config.Feeds {
			fmt.Print("Updating feed: ", feed.Title, "\n")
			err := feed.Update()

			if err != nil {
				fmt.Println("Error in updating RSS feed, see: x/mux/commits.go")
				fmt.Printf("%s\n\n", err)
			} else {
				Notify(p)
			}
			time.Sleep(10 * time.Minute)
		}
		time.Sleep(10 * time.Minute)
	}
}

// Notify subscribed channels to subscribed feeds
func Notify(id int) {
	if !Config.Feeds[id].Items[0].Read {
		resp := ".\n"

		resp += "**" + Config.Feeds[id].Title + ": **" + "\n"
		resp += Config.Feeds[id].Items[0].Date.String() + "\n\n"
		resp += "`" + Config.Feeds[id].Items[0].Title + "`" + "\n"
		resp += "\n" + Config.Feeds[id].Items[0].Link + "\n"
		Config.Feeds[id].Items[0].Read = true
		resp += "\n"
		
		// Loop through subbed chans and post notification message
		for _, v := range Config.Subs {
			if v.SubID == id {
				Session.ChannelMessageSend(v.ChanID, resp)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}


/* Commands for Mux */

// Print last commit for a given feed (by subscribed id int [See: List(...)])
func (m *Mux) Last(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := ".\n"

	id, _ := strconv.Atoi(ctx.Fields[len(ctx.Fields) -1])
	if id >= 0 && id < len(Config.Feeds) {
		resp += "**" + Config.Feeds[id].Title + ": **" + "\n"
		resp += Config.Feeds[id].Items[0].Date.String() + "\n\n"
		resp += "`" + Config.Feeds[id].Items[0].Title + "`" + "\n"
		resp += "\n" + Config.Feeds[id].Items[0].Link + "\n"
		Config.Feeds[id].Items[0].Read = true
		fmt.Println("Last-ing notification: ", Config.Feeds[id].Items[0])
	} else {
		resp += "Denied fetch. Invalid stream id, see: list command"
	}

	resp += "\n"

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

// List all subscribed feeds
func (m *Mux) List(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := "```\n"

	for p, v := range Config.Feeds {
		resp += strconv.Itoa(p) + ": " + v.Title + ", " + v.Link + "\n"
	}

	resp += "```\n"

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

// Subscribe to a feed
func (m *Mux) Add(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
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
	resp += "Added."
	
	resp += "```\n"
	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

// Subscribe current channel to notifications from a given feed id
func (m *Mux) Subscribe(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := "\n"

	id, _ := strconv.Atoi(ctx.Fields[len(ctx.Fields) -1])
	if id >= 0 && id < len(Config.Feeds) {
		var sub Subscription

		// Check if already subscribed
		for _, v := range Config.Subs {
			if v.ChanID == dm.ChannelID && v.SubID == id {
				resp += "Denied subscription. Already subscribed in this channel."
				goto NOSUB
			}
		}
		
		sub.ChanID = dm.ChannelID
		sub.SubID = id
	
		// Might not be thread-safe
		Config.Subs = append(Config.Subs, sub)
		resp += "Subscribed."
		NOSUB:
	} else {
		resp += "Denied subscription. Invalid stream id, see: list command"
	}
	
	resp += "\n"
	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}
