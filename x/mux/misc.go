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

// Dump config to file
func (m *Mux) Dump(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := ""

	err := Config.Write()
	if err != nil {
		resp += "Dump failed. Check logs."
	} else {
		resp += "Ok."
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
