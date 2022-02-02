package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	embed "github.com/Clinet/discordgo-embed"
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

	fmt.Println(string(body))

	channelID := getChannelID("woocommerce")
	_, err = getSession().ChannelMessageSendEmbed(channelID, embed.NewGenericEmbed("Example", "This is an example embed!"))
	if err != nil {
		log.Print(err)
	}

	fmt.Fprintf(w, "thanks\n")
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8090", nil)
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
