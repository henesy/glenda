package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/henesy/glenda/x/mux"
	"log"
	"os"
	"os/signal"
	sc "strconv"
	"syscall"
	"time"
)

const Version = "v0.1.1"

// Session is declared in the global space so it can be easily used
// throughout this program.
// In this use case, there is no error that would be returned.
var Session *discordgo.Session
var token string

// Read in all configuration options from both environment variables and
// command line arguments.
func init() {

	// Discord Authentication Token
	// Have to prefix "Bot [Token Here]" or 401 Forbidden
	token = os.Getenv("DG_TOKEN")
	if token == "" {
		flag.StringVar(&token, "t", "", "Discord Authentication Token (Bot ...)")
	}

	var err error
	Session, err = discordgo.New(token)
	if err != nil {
		log.Fatal("error initiating session")
	}
}

func main() {

	// Declare any variables needed later.
	var err error

	// Print out a fancy logo!
	fmt.Printf(` 
            __
           (  \
      __   \  '\
     (  "-_ \ .-'----._
     '-_  "v"         "-
	"Y'             ".
	 |                |
	 |        o     o |
	 |          .<>.  |
	  \         "Ll"  |
	   |             .'
	   |             |
	   (             /
	  /'\         . \
	  "--^--.__,\_)-'   %-16s\/`+"\n\n", Version)

	// Parse command line arguments
	flag.Parse()

	// Verify a Token was provided
	if Session.Token == "" {
		log.Println("You must provide a Discord authentication token.")
		return
	}

	// Verify the Token is valid and grab user information
	Session.State.User, err = Session.User("@me")
	if err != nil {
		log.Printf("error fetching user information, %s\n", err)
	}

	// Open a websocket connection to Discord
	err = Session.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}

	// Init boot vars
	mux.StartTime = time.Now()

	// Init Mux daemons
	mux.Config.Init(Session)
	mux.RemChan = make(chan mux.Reminder, 5)
	go mux.Reminders()

	Session.UpdateGameStatus(0, "with #cat-v")

	// Wait for a CTRL-C
	log.Printf(`Now running on PID ` + sc.Itoa(os.Getpid()) + `. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	Session.Close()

	// Exit Normally.
}
