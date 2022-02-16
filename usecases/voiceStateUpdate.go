package usecases

import (
	"github.com/bwmarrin/discordgo"
)

const mainChatChannel = "943453314551529494"

func VoiceStateUpdate(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	if m.UserID == s.State.User.ID {
		return
	}
	if m.UserID == "LovelyAlecsmith" {

	}

	s.ChannelMessageSend(mainChatChannel, "Someone joined a voice channel!")
}