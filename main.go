package main

import (
	"log"
	"os"

	"os/signal"
	"syscall"

	"github.com/AlecSmith96/dnd-bot/adapters"
	"github.com/AlecSmith96/dnd-bot/usecases"
	"github.com/bwmarrin/discordgo"
)

func main() {
	config := adapters.GetConfig()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(usecases.MessageCreate)
	dg.AddHandler(usecases.VoiceStateUpdate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
