package mux

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"time"
	"strconv"
	"os/exec"
)

var desc []string = []string{
	"Section (1) for general publicly accessible commands.",
	"Section (2) for library functions, including system calls.",
	"Section (3) for kernel devices (accessed via bind (1)).",
	"Section (4) for file services (accessed via mount).",
	"Section (5) for the Plan 9 file protocol.",
	"Section (6) for file formats.",
	"Section (7) for databases and database access programs.",
	"Section (8) for things related to administering Plan 9."}


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
		url := "http://man.postnix.us/9front/" + section + "/" + page
		
		page, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching URL, see: x/mux/man.go")
			resp += "Error fetching URL. Is man.postnix.us up?"
			goto URL
		}
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
			url := "http://man.postnix.us/9front/" + strconv.Itoa(i) + "/" + page
			
			page, err := http.Get(url)
			defer page.Body.Close()
			body, err := ioutil.ReadAll(page.Body)
			
			if err != nil || strings.Contains(string(body), "The requested document at") {
				continue
			} else {
				if !any {
					any = true
					resp += desc[i-1] + "\n"
				}
				resp += string(url) + "\n"
			}
			time.Sleep(20 * time.Millisecond)
		}
		
		// No matches
		if !any {
			resp += "No matching manual page(s) found."
		}
		
	} else if len(ctx.Fields) == 1 {
		resp += "http://man.postnix.us/9front/\n"
		for i:=1; i < 9; i++ {
			resp += desc[i-1] + "\n"
		}
	} else {
		resp += "No op. Please use the 'man 3 srv' format or 'man srv' format."
	}

	URL:
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
	resp := "\n"

	if len(ctx.Fields) >= 2 {
		// g!lookman vmx
		s := ctx.Fields[1:]
		// Match 10 keys max
		MAXKEYS := 10
		args := make([]string, 0, MAXKEYS)
		for p, v := range s {
			if p < MAXKEYS {
				args = append(args, v)
			} else {
				break
			}
		}
	
		// rc or get out
		//fmt.Println("Running: ", args)
		out, err := exec.Command("./x/mux/lookman/lookman", args...).Output()
		if err != nil {
			resp += "No lookman script found."
			goto END
		} else {
			if len(out) < 2 {
				resp += "No matching manual page(s) found."
				goto END
			}
			
			top := true
			str := string(out)
			lines := strings.Split(str, "\n")
			for p, l := range lines {
				if p == len(lines) -1 {
					break
				}
				fields := strings.Split(l, " ")
				i := fields[1]
				page := fields[2]
				
				url := "http://man.postnix.us/9front/" + i + "/" + page + " # " + fields[4] + "\n"
	
				if top {
					top = false
					first, _ := strconv.Atoi(i)
					resp += desc[first-1] + "\n"
				}
				
				resp += url
			}
		}
	} else {
		resp += "Usage: lookman key ...\n"
	}
	
	if len(resp) > 2000 {
		resp = "\nError: Lookman output exceeded Discord™ 2000 character limit."
	}

	END:
	resp += "\n"
	_, _ = ds.ChannelMessageSend(dm.ChannelID, resp)
}

// Call lookman on a given query
func (m *Mux) Mkindex(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := "\n"

	out, err := exec.Command("./gendex").Output()
	if err != nil {
		fmt.Println(err)
		resp += "No mkindex script found."
		goto END
	} else {
		fmt.Println("Generating out: ", out)
		resp += "Ok."
	}

	END: 
	resp += "\n"
	ds.ChannelMessageSend(dm.ChannelID, resp)
}
