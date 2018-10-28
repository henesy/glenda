package mux

import (
	"github.com/bwmarrin/discordgo"
	"github.com/SlyMarbo/rss"
	"fmt"
	"time"
	"strings"
	"strconv"
	"html"
)

/* internal functions */

// Listen on the rss feed
func Listener() {
	for {
		for p, _ := range Feeds.Feeds {
			//fmt.Print("Updating feed: ", Feeds.Feeds[p].Feed.Title, "\n")
			// maybe only do at init step?
			str := Feeds.Feeds[p].Feed.UpdateURL
			feed, err := rss.Fetch(str)
			if feed != nil {
				Feeds.Feeds[p].Feed = *feed
			} else {
				//fmt.Println("Got a nil pointer for feed in commits for ", str, " as ", err)
			}
			//err := Feeds.Feeds[p].Feed.Update()

			if err != nil {
				fmt.Println("Error in updating RSS feed, see: x/mux/commits.go")
				fmt.Printf("%s\n\n", err)
			} else {
				additem := true
				for j, _ := range Feeds.Feeds[p].Recent {
					if Feeds.Feeds[p].Recent[j] == Feeds.Feeds[p].Feed.Items[0].Title {
						//fmt.Println("Checking Recent ", j, " against ", Feeds.Feeds[p].Feed.Items[0].Title)
						additem = false
						break
					}
				}
				if additem {
					// x y z → x y z 0 → y z 0
					fmt.Println("Updating: ", Feeds.Feeds[p].Feed.Title)
					Feeds.Feeds[p].Recent = append(Feeds.Feeds[p].Recent, Feeds.Feeds[p].Feed.Items[0].Title)
					Feeds.Feeds[p].Recent = Feeds.Feeds[p].Recent[1:]
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
	resp += "**" + Feeds.Feeds[id].Feed.Title + ": **" + "\n"
	resp += Feeds.Feeds[id].Feed.Items[0].Date.String() + "\n\n"
	// If a 9front feed, extract the user ☺
	if strings.Contains(Feeds.Feeds[id].Feed.Items[0].Link, "http://code.9front.org/hg/") {
		lines := strings.Split(Feeds.Feeds[id].Feed.Items[0].Summary, "\n")
		for i, v := range lines {
			if strings.Contains(v, "<th style=\"text-align:left;vertical-align:top;\">user</th>") {
				line := html.UnescapeString((lines[i+1])[6:len(lines[i+1])-5])
				resp += line + "\n\n"
				break
			}
		}
	}
	resp += "`" + Feeds.Feeds[id].Feed.Items[0].Title + "`" + "\n"
	resp += "\n" + Feeds.Feeds[id].Feed.Items[0].Link + "\n"
	Feeds.Feeds[id].Feed.Items[0].Read = true
	resp += "\n"
	
	// Loop through subbed chans and post notification message
	fmt.Println("Looping through subs to notify...")
	for _, v := range Feeds.Subs {
		if v.SubID == id {
			Session.ChannelMessageSend(v.ChanID, resp)
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("No new notifys for ", Feeds.Feeds[id].Feed.UpdateURL)
	fmt.Println(Feeds.Feeds[id].Feed.Items[0])
	fmt.Println(Feeds.Feeds[id].Feed.Items[len(Feeds.Feeds[id].Feed.Items)-1])
}


/* Commands for Mux */

// Print last commit for a given feed (by subscribed id int [See: List(...)])
func (m *Mux) Last(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := ".\n"

	id, _ := strconv.Atoi(ctx.Fields[len(ctx.Fields) -1])
	if id >= 0 && id < len(Feeds.Feeds) {
		resp += "**" + Feeds.Feeds[id].Feed.Title + ": **" + "\n"
		resp += Feeds.Feeds[id].Feed.Items[0].Date.String() + "\n\n"
		// If a 9front feed, extract the user ☺
		if strings.Contains(Feeds.Feeds[id].Feed.Items[0].Link, "http://code.9front.org/hg/") {
			lines := strings.Split(Feeds.Feeds[id].Feed.Items[0].Summary, "\n")
			for i, v := range lines {
				if strings.Contains(v, "<th style=\"text-align:left;vertical-align:top;\">user</th>") {
					line := html.UnescapeString((lines[i+1])[6:len(lines[i+1])-5])
					resp += line + "\n\n"
					break
				}
			}
		}
		resp += "`" + Feeds.Feeds[id].Feed.Items[0].Title + "`" + "\n"
		//resp += "\n" + Feeds.Feeds[id].Feed.Items[0].Summary + "\n"
		resp += "\n" + Feeds.Feeds[id].Feed.Items[0].Link + "\n"
		Feeds.Feeds[id].Feed.Items[0].Read = true
		fmt.Println("Last-ing notification: ", Feeds.Feeds[id].Feed.Items[0])
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

	for p, v := range Feeds.Feeds {
		resp += strconv.Itoa(p) + ": " + v.Feed.Title + ", " + v.Feed.Link + "\n"
	}

	resp += "```\n"

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

// Add a given feed to be tracked within glenda
func (m *Mux) Add(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	//http://code.9front.org/hg/plan9front/rss-log
	resp := "```\n"
	// URL to feed should be last item
	url := ctx.Fields[len(ctx.Fields) -1]
	fmt.Println("Proposed addition for: ", url)
	
	for _, v := range Feeds.Feeds {
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
	Feeds.Feeds = append(Feeds.Feeds, tmpFeed)
	resp += "Added."
	
	resp += "```\n"
	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

// Stop tracking a given feed
func (m *Mux) Remove(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	//http://code.9front.org/hg/plan9front/rss-log
	resp := "```\n"
	// ID to feed should be last item
	ID, err := strconv.Atoi(ctx.Fields[len(ctx.Fields) -1])
	fmt.Println("Proposed removal for: ", ID)
	
	if ID < 0 || ID > len(Feeds.Feeds) || err != nil {
		resp += "Denied! Feed does not exist!"
		resp += "```\n"
		goto REND;
	}

	

	// Might not be thread safe
	Feeds.Feeds = append(Feeds.Feeds[:ID], Feeds.Feeds[ID+1:]...)
	resp += "Removed."
	
	REND:
	resp += "```\n"
	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

// Subscribe current channel to notifications from a given feed id
func (m *Mux) Subscribe(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := "\n"

	id, _ := strconv.Atoi(ctx.Fields[len(ctx.Fields) -1])
	if id >= 0 && id < len(Feeds.Feeds) {
		var sub Subscription

		// Check if already subscribed
		for _, v := range Feeds.Subs {
			if v.ChanID == dm.ChannelID && v.SubID == id {
				resp += "Denied subscription. Already subscribed in this channel."
				goto NOSUB
			}
		}
		
		sub.ChanID = dm.ChannelID
		sub.SubID = id
	
		// Might not be thread-safe
		Feeds.Subs = append(Feeds.Subs, sub)
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
	if id >= 0 && id < len(Feeds.Feeds) {
		// Check if subscribed
		removed := false
		for i, v := range Feeds.Subs {
			if v.ChanID == dm.ChannelID && v.SubID == id {
				removed = true
				Feeds.Subs = Feeds.Subs[:i+copy(Feeds.Subs[i:], Feeds.Subs[i+1:])]
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
