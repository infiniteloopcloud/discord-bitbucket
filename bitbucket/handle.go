package bitbucket

import (
	"encoding/json"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

func Handle(eventType string, body []byte) (*discordgo.MessageEmbed, error) {
	switch eventType {
	case "repo:push":
		return handlePush(body)
	}
	return nil, nil
}

func handlePush(body []byte) (*discordgo.MessageEmbed, error) {
	var push RepoPushEvent
	err := json.Unmarshal(body, &push)
	if err != nil {
		return nil, err
	}
	return embed.NewEmbed().SetTitle(push.Repository.Name + " - Push").SetDescription(push.Actor.DisplayName + "pushed").MessageEmbed, nil
}
