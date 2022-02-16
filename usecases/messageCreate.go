package usecases

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/Clinet/discordgo-embed"
)

const commandWord = "dnd-bot"

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
//
// It is called whenever a message is created but only when it's sent through a
// server as we did not request IntentsDirectMessages.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	firstWord := strings.Split(m.Content, " ")[0]
	words := strings.Split(m.Content, " ")

	// if message from own bot or not a bot command
	if m.Author.ID == s.State.User.ID ||  firstWord != commandWord {
		return
	}

	switch words[1] {
	case "ping": 
		log.Print("ping command entered")
		messagePong(s, m)
	default:
		log.Print("Unrecognised command entered")
		content := fmt.Sprintf("Unrecognised command: %s", words[1])
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewErrorEmbed("Unrecognised", content))
	}
}

func messagePong(s *discordgo.Session, m *discordgo.MessageCreate) {
	// We create the private channel with the user who sent the message.
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
		// If an error occurred, we failed to send the message.
		//
		// It may occur either when we do not share a server with the
		// user (highly unlikely as we just received a message) or
		// the user disabled DM in their settings (more likely).
		log.Println("error sending DM message:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Failed to send you a DM. "+
				"Did you disable DM in your privacy settings?",
		)
	}
}
