package usecases

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"github.com/AlecSmith96/dnd-bot/adapters"
)

const mainChatChannel = "943453314551529494"

func VoiceStateUpdate(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	config := adapters.GetConfig()

	if m.UserID == s.State.User.ID {
		return
	}

	// don't send message if they are leaving the channel
	if m.BeforeUpdate != nil {
		return
	}
	
	// send request to get user 
	client := http.Client{}
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://discord.com/api/guilds/%s", m.GuildID), nil)
	headerContent := fmt.Sprintf("Bot %s", config.Token)
	req.Header = http.Header{
		"Authorization": []string{headerContent},
	}
	resp , err := client.Do(req)
	if err != nil {
		log.Printf("unable to send GET request: %v", err)

	}


	fmt.Print(resp.Body)

	// read json http response
	jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Print("unable to read json data")
	}

	fmt.Print(jsonDataFromHttp)

	var guild discordgo.Guild
	err = json.Unmarshal([]byte(jsonDataFromHttp), &guild) // here!

	if err != nil {
			panic(err)
	}

	fmt.Print(guild)

	s.ChannelMessageSendEmbed(mainChatChannel, embed.NewGenericEmbed("Alert", "Someone joined a voice channel!"))
}