package mux

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"os/exec"
	"fmt"
)

// Display fortunes from various 9front /lib/ files
func (m *Mux) Fortunes(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := "```\n"

	// Read relevant file
	path := "/usr/share/mirror/plan9front/lib/"

	f := ""
	if strings.Contains(dm.Content, "theo") {
		f = "theo"
	} else if strings.Contains(dm.Content, "troll") {
		f = "troll"
	} else if strings.Contains(dm.Content, "rsc") {
		f = "rsc"
	} else if strings.Contains(dm.Content, "terry") {
		f = "terry"
	} else if strings.Contains(dm.Content, "rob") {
		f = "rob"
	} else if strings.Contains(dm.Content, "ken") {
		f = "ken"
	} else if strings.Contains(dm.Content, "henesy") {
		f = "henesy"
	} else if strings.Contains(dm.Content, "davros") {
		f = "davros"
	} else {
		f = "../sys/games/lib/fortunes"
	}

	// Linux solution
	out, err := exec.Command("fortune", path + f).Output()
	if err != nil {
		fmt.Println("Error calling fortune(1), see: x/mux/fortunes.go")
		fmt.Println("%s\n", err)
	}
	resp += string(out)

	resp += "```\n"

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

// Display bullshit from 9front
func (m *Mux) Bullshit(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := "```\n"
	
	out, err := exec.Command("./x/mux/misc/bullshit").Output()
	if err != nil {
		resp += "No bullshit script found."
		goto END
	} else {
		resp += string(out)
	}
	
	END:
	resp += "```\n"
	ds.ChannelMessageSend(dm.ChannelID, resp)
	return
}
