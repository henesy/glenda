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
		for p, _ := range Config.Feeds {
			//fmt.Print("Updating feed: ", Config.Feeds[p].Feed.Title, "\n")
			// maybe only do at init step?
			str := Config.Feeds[p].Feed.UpdateURL
			feed, err := rss.Fetch(str)
			Config.Feeds[p].Feed = *feed
			//err := Config.Feeds[p].Feed.Update()

			if err != nil {
				fmt.Println("Error in updating RSS feed, see: x/mux/commits.go")
				fmt.Printf("%s\n\n", err)
			} else {
				additem := true
				for j, _ := range Config.Feeds[p].Recent {
					if Config.Feeds[p].Recent[j] == Config.Feeds[p].Feed.Items[0].Title {
						//fmt.Println("Checking Recent ", j, " against ", Config.Feeds[p].Feed.Items[0].Title)
						additem = false
						break
					}
				}
				if additem {
					// x y z → x y z 0 → y z 0
					fmt.Println("Updating: ", Config.Feeds[p].Feed.Title)
					Config.Feeds[p].Recent = append(Config.Feeds[p].Recent, Config.Feeds[p].Feed.Items[0].Title)
					Config.Feeds[p].Recent = Config.Feeds[p].Recent[1:]
					Notify(p)
					break
				}
			}
			time.Sleep(2 * time.Minute)
		}
		time.Sleep(10 * time.Minute)
		// Dump config to file regularly
		Config.Write()
	}
}

// Notify subscribed channels to subscribed feeds
func Notify(id int) {
	resp := ".\n"
	resp += "**" + Config.Feeds[id].Feed.Title + ": **" + "\n"
	resp += Config.Feeds[id].Feed.Items[0].Date.String() + "\n\n"
	resp += "`" + Config.Feeds[id].Feed.Items[0].Title + "`" + "\n"
	resp += "\n" + Config.Feeds[id].Feed.Items[0].Link + "\n"
	Config.Feeds[id].Feed.Items[0].Read = true
	resp += "\n"
	
	// Loop through subbed chans and post notification message
	fmt.Println("Looping through subs to notify...")
	for _, v := range Config.Subs {
		if v.SubID == id {
			Session.ChannelMessageSend(v.ChanID, resp)
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("No new notifys for ", Config.Feeds[id].Feed.UpdateURL)
	fmt.Println(Config.Feeds[id].Feed.Items[0])
	fmt.Println(Config.Feeds[id].Feed.Items[len(Config.Feeds[id].Feed.Items)-1])
}


/* Commands for Mux */

// Print last commit for a given feed (by subscribed id int [See: List(...)])
func (m *Mux) Last(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := ".\n"

	id, _ := strconv.Atoi(ctx.Fields[len(ctx.Fields) -1])
	if id >= 0 && id < len(Config.Feeds) {
		resp += "**" + Config.Feeds[id].Feed.Title + ": **" + "\n"
		resp += Config.Feeds[id].Feed.Items[0].Date.String() + "\n\n"
		resp += "`" + Config.Feeds[id].Feed.Items[0].Title + "`" + "\n"
		resp += "\n" + Config.Feeds[id].Feed.Items[0].Link + "\n"
		Config.Feeds[id].Feed.Items[0].Read = true
		fmt.Println("Last-ing notification: ", Config.Feeds[id].Feed.Items[0])
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
		resp += strconv.Itoa(p) + ": " + v.Feed.Title + ", " + v.Feed.Link + "\n"
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
		// this is bad matching, can't have two bitbucket url's?
		if strings.Contains(url, v.Feed.UpdateURL) {
			//fmt.Println(url)
			//fmt.Println(v.Link)
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
	var tmpFeed Feed
	tmpFeed.Feed = *feed
	// Maybe make the size here a Config variable
	tmpFeed.Recent = make([]string, 3)
	Config.Feeds = append(Config.Feeds, tmpFeed)
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

// a = a[:i+copy(a[i:], a[i+1:])]
// Unsubscribe current channel to notifications from a given feed id
func (m *Mux) Unsubscribe(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := "\n"

	id, _ := strconv.Atoi(ctx.Fields[len(ctx.Fields) -1])
	if id >= 0 && id < len(Config.Feeds) {
		// Check if subscribed
		removed := false
		for i, v := range Config.Subs {
			if v.ChanID == dm.ChannelID && v.SubID == id {
				removed = true
				Config.Subs = Config.Subs[:i+copy(Config.Subs[i:], Config.Subs[i+1:])]
			}
		}
		
		if(!removed) {
			resp += "Denied unsubscription. Not subscribed in this channel."
		} else {
			resp += "Unsubscribed."
		}
	} else {
		resp += "Denied unsubscription. Invalid stream id, see: list command"
	}
	
	resp += "\n"
	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}
