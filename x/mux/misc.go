package mux

import (
	"github.com/bwmarrin/discordgo"
	"time"
	"strings"
)

// Chan for communicating to Glenda
var GlendaChan chan string

// Chan for communicating to Mux
var MuxChan chan string

// Chan for communicating with Dump
var dumpChan chan string

// Stores the time that the bot started this boot
var StartTime time.Time


// Multiplex internal channels, initialized once
func CommMux() {
	MuxChan = make(chan string, 5)
	GlendaChan = make(chan string, 5)
	dumpChan = make(chan string)

	// Listen for signals till death do us part
	for {
		select {
		default:
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// Dump configs to file
func dump() string {
	resp := ""

	err := Config.Write()
	if err != nil {
		resp += "Dump config failed. Check logs.\n"
	} else {
		resp += "Ok."
	}
	resp += "\n"
	return resp
}

// Return the current uptime
func (m *Mux) Uptime(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := ""
	
	resp += time.Now().Sub(StartTime).String()

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

// Check authorization as owner
func authorized(dm *discordgo.Message) bool {
	user := dm.Author
	id := user.ID
	
	if id == Config.Db["owner"] {
		return true
	}
	
	return false
}

// Dump config to file
func (m *Mux) Dump(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	if !authorized(dm) {
		ds.ChannelMessageSend(dm.ChannelID, "Only the bot owner can do that.")
		return
	}

	resp := dump()

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

// Rewrite griddisk paths as links
// Given /n/griddisk/foo, /foo, foo ⇒ …/incoming/foo
func (m *Mux) GridLink(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := ""
	base := "http://wiki.9gridchan.org/incoming"
	root := "/n/griddisk"
	
	resp += base
	
	path := ctx.Fields[1]
	
	switch {
	case strings.HasPrefix(path, root):
		// Full path
		sub := strings.TrimPrefix(path, root)
		resp += sub
	
	case path[0] == '/':
		// Root only
		resp += path
	
	default:
		// Any other case
		resp += "/" + path
	}
	
	resp += "\n"

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

// Beer, dude.
func (m *Mux) Beer(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := ""
	resp += ":beer:"
	resp += "\n"

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

// Whiskey, dude.
func (m *Mux) Whiskey(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := ""
	resp += ":tumbler_glass:"
	resp += "\n"

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

// Wine, dude.
func (m *Mux) Wine(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := ""
	resp += ":wine_glass:"
	resp += "\n"

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}
