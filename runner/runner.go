package runner

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/infiniteloopcloud/discord-bitbucket/bitbucket"
	"github.com/infiniteloopcloud/discord-bitbucket/env"
)

var session *discordgo.Session
var channelsCache map[string]string

func webhookHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
	}

	event := req.Header.Get("X-Event-Key")
	channel, message, err := bitbucket.Handle(event, body)
	if err != nil {
		log.Printf("[ERROR] [%s] %s", event, err.Error())
		return
	}
	channelID := getChannelID(channel)
	if channelID == "" {
		channelID = getChannelID("unknown")
	}
	if channelID != "" && message != nil {
		_, err = getSession().ChannelMessageSendEmbed(channelID, message)
		if err != nil {
			log.Printf("[ERROR] [%s] %s", event, err.Error())
		}
	}

	fmt.Fprintf(w, "ACK")
}

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ACK")
}

func getChannelID(name string) string {
	if channelsCache == nil {
		channelsCache = make(map[string]string)
	}
	if id, ok := channelsCache[name]; ok {
		return id
	} else {
		channels, err := getSession().GuildChannels(env.Configuration().BotGuild)
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
		session, err = discordgo.New("Bot " + env.Configuration().BotToken)
		if err != nil {
			log.Printf("[ERROR] %s", err.Error())
		}
	}
	return session
}

func Run() {
	env.Configuration().Dump()

	address := ":8080"
	if a := env.Configuration().Address; a != "" {
		address = a
	}

	http.HandleFunc("/bitbucket/webhooks", webhookHandler)
	http.HandleFunc("/bitbucket/hc", healthCheckHandler)
	log.Printf("Server listening on %s", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Printf("[ERROR] %s", err.Error())
	}
}
