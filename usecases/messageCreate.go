package usecases

import (
	"fmt"
	"log"
	"strings"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

const commandWord = "dbot"

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
//
// It is called whenever a message is created but only when it's sent through a
// server as we did not request IntentsDirectMessages.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	firstWord := strings.Split(m.Content, " ")[0]
	words := strings.Split(m.Content, " ")

	// if message from own bot or not a bot command
	if m.Author.ID == s.State.User.ID || firstWord != commandWord {
		return
	}

	log.Print(m.Author.ID)

	switch words[1] {
	case "ping":
		log.Print("ping command entered")
		messagePong(s, m)
	case "create":
		handleCreate(s, m, words)
	case "next":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("DnD Bot", "Your next session is scheduled for <session_timestamp>"))
	default:
		log.Print("Unrecognised command entered")
		content := fmt.Sprintf("Unrecognised command: %s", words[1])
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewErrorEmbed("Unrecognised", content))
	}
}

// handleCreate checks to see what the user wants to create and creates it
func handleCreate(s *discordgo.Session, m *discordgo.MessageCreate, words []string) {
	switch words[2] {
	case "session":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("New Session Created", "A new session has been created for <session_timestamp>"))
	default:
		content := fmt.Sprintf("Unable to create record for: %s", words[2])
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewErrorEmbed("Unrecognised", content))
	}
}

// messagePong send a private message to the author of the message.
func messagePong(s *discordgo.Session, m *discordgo.MessageCreate) {
	// create the private channel with the user who sent the message
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		log.Println("error creating channel:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Something went wrong while sending the DM!",
		)
		return
	}
	// Then we send the message through the channel we created.
	_, err = s.ChannelMessageSend(channel.ID, fmt.Sprintf("Hi "+m.Author.Username+": Pong!"))

	if err != nil {
		log.Println("error sending DM message:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Failed to send you a DM "+
				m.Author.Username+
				". Did you disable DM in your privacy settings?",
		)
	}
}
