package mux

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"os/exec"
	"fmt"
)

// Display conjured fortunes from the oracle
// g!oracle papers factotum is
func (m *Mux) Oracle(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := "```\n"

	// Read relevant file
	path := Config.Db["oracledbs"]

	fields := strings.Fields(dm.Content)
	fields = fields[1:]
	f := "papers"
	prompt := []string{}
	if len(fields) > 1 {
		prompt = fields[1:]
	}

	if strings.Contains(dm.Content, "papers") {
		f = "papers"
	} else if strings.Contains(dm.Content, "fqa") {
		f = "fqa"
	} else if strings.Contains(dm.Content, "pooh") {
		f = "pooh"
	} else if strings.Contains(dm.Content, "faust") {
		f = "faust"
	} else if strings.Contains(dm.Content, "fortunes") {
		f = "fortunes"
	}

	// Linux solution
	var args []string
	args = append(args, "-db", path + f, "-len", "30")
	args = append(args, prompt...)
	//fmt.Println(args)
	out, err := exec.Command("oracle", args...).Output()
	if err != nil {
		fmt.Println("Error calling oracle(1), see: x/mux/oracle.go")
		fmt.Println("%s\n", err)
	}
	resp += string(out)

	resp += "```\n"

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}


