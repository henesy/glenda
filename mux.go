package main

// This file adds the Disgord message route multiplexer, aka "command router".
// to the Disgord bot. This is an optional addition however it is included
// by default to demonstrate how to extend the Disgord bot.

import "bitbucket.org/henesy/glenda/x/mux"

// Router is registered as a global variable to allow easy access to the
// multiplexer throughout the bot.
var Router = mux.New()

func init() {
	// Register the mux OnMessageCreate handler that listens for and processes
	// all messages received.
	Session.AddHandler(Router.OnMessageCreate)

	// Register the build-in help command.
	Router.Route("help", "Display this message.", Router.Help)

	Router.Route("fortune", "Display fortunes. Bonus files are (theo troll rsc bullshit terry rob).", Router.Fortunes)
	
	Router.Route("subscribe", "Subscribe to an RSS feed.", Router.Subscribe)
	
	Router.Route("list", "List RSS feeds Glenda is subscribed to by id.", Router.List)
	
	Router.Route("last", "Show the last commit for a given RSS feed by id.", Router.Last)

	Router.Route("dump", "Dump config to file.", Router.Dump)

	Router.Route("about", "General information about the bot.", Router.About)
	

}
