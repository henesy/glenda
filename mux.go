package main

// This file adds the Disgord message route multiplexer, aka "command router".
// to the Disgord bot. This is an optional addition however it is included
// by default to demonstrate how to extend the Disgord bot.

import (
	"github.com/henesy/glenda/x/mux"
)

// Router is registered as a global variable to allow easy access to the
// multiplexer throughout the bot.
var Router = mux.New()

func init() {
	// Register the mux OnMessageCreate handler that listens for and processes
	// all messages received.
	Session.AddHandler(Router.OnMessageCreate)

	// Register the build-in help command.
	Router.Route("help", "Display this message.", Router.Help)

	Router.Route("about", "General information about the bot.", Router.About)

	Router.Route("newuser", "Information for new users of Plan 9.", Router.Newuser)

	Router.Route("remindme", "Set a reminder for a given time interval ([int][hmsd] [reminder]).", Router.RemindMe)

	Router.Route("man", "Link a given Plan 9 manual page (From 9front for now).", Router.Man)

	Router.Route("lookman", "Search for a string in the Plan 9 manual pages (From 9front for now).", Router.Lookman)

	Router.Route("sig", "Search for function definitions in the Plan 9 manual pages (From 9front for now).", Router.Sig)

	Router.Route("fortune", "Display fortunes. Bonus files are (theo troll rsc terry rob ken henesy davros).", Router.Fortunes)

	Router.Route("bullshit", "Print logical statements with sound grounding.", Router.Bullshit)

	Router.Route("add", "Track an RSS feed.", Router.Add)

	Router.Route("remove", "Remove a tracked RSS feed by ID.", Router.Remove)

	Router.Route("list", "List RSS feeds Glenda is subscribed to by id.", Router.List)

	Router.Route("last", "Show the last commit for a given RSS feed by id.", Router.Last)

	Router.Route("dump", "Dump config to file.", Router.Dump)

	Router.Route("mkindex", "Generate index for lookman from mirror.", Router.Mkindex)

	Router.Route("subscribe", "Subscribe current channel to a given feed id.", Router.Subscribe)

	Router.Route("unsubscribe", "Unsubscribe current channel to a given feed id.", Router.Unsubscribe)

	Router.Route("beer", ":beer:", Router.Beer)

	Router.Route("whiskey", ":tumbler_glass:", Router.Whiskey)

	Router.Route("wine", ":wine_glass:", Router.Wine)

	Router.Route("uptime", "Current bot uptime", Router.Uptime)
	
	Router.Route("gridlink", "Convert griddisk paths to 'incoming' URL's", Router.GridLink)

	Router.Route("roll", "Roll dice in the form XdY", Router.Roll)
}
