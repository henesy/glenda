package mux

import (
	"github.com/bwmarrin/discordgo"
)

// Generate a 1 time use invite
func (m *Mux) Invite(ds *discordgo.Session, dm *discordgo.Message, ctx *Context) {
	invite := discordgo.Invite{MaxUses: 1}
	st, err := ds.ChannelInviteCreate(dm.ChannelID, invite)
	if err != nil {
		ds.ChannelMessageSend(dm.ChannelID, "Error generating invite: "+err.Error())
		return
	}

	ds.ChannelMessageSend(dm.ChannelID, "https://discord.gg/"+st.Code+"\n")
}
