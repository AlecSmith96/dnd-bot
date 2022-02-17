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

	// if the bot has joined a channel, return
	if m.UserID == s.State.User.ID {
		return
	}

	// don't send message if they are leaving the channel
	if m.BeforeUpdate != nil {
		return
	}
	
	// send request to get user 
	client := http.Client{}
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://discord.com/api/guilds/%s/members/%s", m.GuildID, m.UserID), nil)
	headerContent := fmt.Sprintf("Bot %s", config.Token)
	req.Header = http.Header{
		"Authorization": []string{headerContent},
	}
	resp , err := client.Do(req)
	if err != nil {
		log.Printf("unable to send GET request: %v", err)

	}

	// read json http response
	jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("unable to read json data")
	}

	var member discordgo.Member
	err = json.Unmarshal([]byte(jsonDataFromHttp), &member)
	if err != nil {
		log.Printf("Error unmarshalling json response: %v", err)
	}

	if len(member.Roles) == 0 {
		return
	}

	dmRoleID := getDMRoleID(client, m.GuildID, config.Token)

	// check that user is the DM before sending message
	for _, role := range member.Roles {
		if role == dmRoleID {
			s.ChannelMessageSendEmbed(mainChatChannel, embed.NewGenericEmbed("Adventurers!", "The DM has entered the Inn!"))
		}
	}
}

// getDMRoleID gets the specific role ID for the 'DM' role.
func getDMRoleID(client http.Client, guildID string, token string) string {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://discord.com/api/guilds/%s/roles", guildID), nil)
	headerContent := fmt.Sprintf("Bot %s", token)
	req.Header = http.Header{
		"Authorization": []string{headerContent},
	}
	resp , err := client.Do(req)
	if err != nil {
		log.Printf("unable to send GET request: %v", err)

	}

	// read json http response
	jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("unable to read json data")
	}

	var roles discordgo.Roles
	err = json.Unmarshal([]byte(jsonDataFromHttp), &roles)
	if err != nil {
		log.Printf("Error unmarshalling json response: %v", err)
	}

	for _, role := range roles {
		if role.Name == "DM" {
			return role.ID
		}
	}

	return ""
}