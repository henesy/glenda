package mux

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"time"
	"strconv"
)


// Fetch a man page
func (m *Mux) Man(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := "\n"
	var err error 
	
	page := ""
	section := ""
	if len(ctx.Fields) == 3 {
		// `man 1 bc` format
		section += ctx.Fields[len(ctx.Fields)-2]
		page += ctx.Fields[len(ctx.Fields)-1]
	
		// Should allow selection of manual categories beyond 9front
		url := "http://man.cat-v.org/9front/" + section + "/" + page
		
		page, err := http.Get(url)
		defer page.Body.Close()
		body, err := ioutil.ReadAll(page.Body)
		
		if err != nil || strings.Contains(string(body), "The requested document at") {
			fmt.Println("Error finding man page, see: x/mux/man.go")
			fmt.Println("%s\n", err)
			resp += "Invalid manual page requested."
		} else {
			resp += string(url)
		}
		
	} else if len(ctx.Fields) == 2 {
		// `man srv` format
		page += ctx.Fields[len(ctx.Fields)-1]
		any := false
		
		for i:=1; i < 9; i++ {
			// Should allow selection of manual categories beyond 9front
			url := "http://man.cat-v.org/9front/" + strconv.Itoa(i) + "/" + page
			
			page, err := http.Get(url)
			defer page.Body.Close()
			body, err := ioutil.ReadAll(page.Body)
			
			if err != nil || strings.Contains(string(body), "The requested document at") {
				continue
			} else {
				if !any {
					any = true
					resp += ".\n"
				}
				resp += string(url) + "\n"
			}
			time.Sleep(1 * time.Second)
		}
		
		// No matches
		if !any {
			resp += "No matching manual page(s) found."
		}
		
	} else {
		resp += "No op. Please use the 'man 3 srv' format or 'man srv' format."
	}

	resp += "\n"

	_, err = ds.ChannelMessageSend(dm.ChannelID, resp)
	if err != nil {
		fmt.Println("Error sending man(1) output as message, see: x/mux/man.go")
		fmt.Println(err)
	}

	return
}

// Fetch a summary of a man page
func (m *Mux) Sum(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	
}

// Call lookman on a given query
func (m *Mux) Lookman(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	
}
