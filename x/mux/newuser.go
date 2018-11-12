package mux

import (
	"github.com/bwmarrin/discordgo"
)

// Help for newbies. :D
func (m *Mux) Newuser(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	resp := ""
	resp += `Welcome to 9fans! Here's an introduction and resources for Plan 9 and 9front.
In this community, you are expected to read the docs and man pages.
However, don't be afraid to ask for help and clarification!
Install 9front: http://9front.org/ http://fqa.9front.org/
Beginner's guides: http://lsub.org/who/nemo/9.intro.pdf http://fqa.9front.org/fqa8.html
Read the man pages: http://man.9front.org/ (also use 'man [page name]' or 'lookman [keyword]')
Read the docs: http://doc.cat-v.org/plan_9/4th_edition/papers/ (also found in /sys/doc, use 'page [file]')
Wiki for additional help: https://9p.io/wiki/plan9/plan_9_wiki/ (some info may be outdated)"
Miscellaneous Plan 9 websites: http://blog.postnix.us/ https://code.9front.org/hg/
http://9gridchan.org/ http://felloff.net/usr/cinap_lenrek/
	`
	resp += "\n"

	ds.ChannelMessageSend(dm.ChannelID, resp)

	return
}

