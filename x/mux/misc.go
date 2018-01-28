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
		case r := <-MuxChan:
			if r == "dump ok" {
				dumpChan <- "Ok."
			} else if r == "dump failed" {
				dumpChan <- "Dump Failed. Check Logs."
			}
		case <-dumpChan:
			GlendaChan <- "dump"
		default:
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// Dump config to file
func (m *Mux) Dump(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := ""

	dumpChan <- ""
	resp += <-dumpChan

	resp += "\n"

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}



