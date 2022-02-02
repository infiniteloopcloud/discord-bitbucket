package main

import (
	"fmt"
	"github.com/infiniteloopcloud/discord-bitbucket/bitbucket"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

const (
	Token = "TOKEN"
	// Get Guild ID for security reasons
	GuildID = "GUILD_ID"
)

var session *discordgo.Session
var channelsCache map[string]string
var guildID = "938346153509015552"

func hello(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print(err)
	}

	channelID := getChannelID("woocommerce")
	message, err := bitbucket.Handle(req.Header.Get("X-Event-Key"), body)
	if err != nil {
		log.Print(err)
	}
	_, err = getSession().ChannelMessageSendEmbed(channelID, message)
	if err != nil {
		log.Print(err)
	}

	fmt.Fprintf(w, "thanks\n")
}

func main() {
	http.HandleFunc("/webhooks", hello)
	http.ListenAndServe(":8000", nil)
}

func getChannelID(name string) string {
	if channelsCache == nil {
		channelsCache = make(map[string]string)
	}
	if id, ok := channelsCache[name]; ok {
		return id
	} else {
		channels, err := getSession().GuildChannels(guildID)
		if err != nil {
			log.Print(err)
		}
		for _, channel := range channels {
			if name == channel.Name {
				channelsCache[channel.Name] = channel.ID
				return channel.ID
			}
		}
	}
	return ""
}

func getSession() *discordgo.Session {
	if session == nil {
		var err error
		session, err = discordgo.New("Bot " + os.Getenv(Token))
		if err != nil {
			panic(err)
		}
	}
	return session
}
