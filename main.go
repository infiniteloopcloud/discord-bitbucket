package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/infiniteloopcloud/discord-bitbucket/bitbucket"
)

const (
	Token = "TOKEN"
	// Get Guild ID for security reasons
	GuildID = "GUILD_ID"

	Address = "ADDRESS"
)

var session *discordgo.Session
var channelsCache map[string]string

//var guildID = "938346153509015552"

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
	address := ":8080"
	if a := os.Getenv(Address); a != "" {
		address = a
	}

	http.HandleFunc("/webhooks", hello)
	log.Printf("Server listening on %s", address)
	http.ListenAndServe(address, nil)
}

func getChannelID(name string) string {
	if channelsCache == nil {
		channelsCache = make(map[string]string)
	}
	if id, ok := channelsCache[name]; ok {
		return id
	} else {
		channels, err := getSession().GuildChannels(os.Getenv(GuildID))
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
