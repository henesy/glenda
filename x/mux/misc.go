package mux

import (
	"github.com/bwmarrin/discordgo"
	"time"
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
	err = Feeds.Write()
	if err != nil {
		resp += "Dump feeds failed. Check logs.\n"
	} else {
		resp += "Ok."
	}
	err = RemindersWrite()
	if err != nil {
		resp += "Dump reminders failed. Check logs.\n"
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

// Dump config to file
func (m *Mux) Dump(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := dump()

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
