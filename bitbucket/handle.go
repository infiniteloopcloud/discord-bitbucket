package bitbucket

import (
	"encoding/json"
	"fmt"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

const (
	success = 0x90EE90
	failure = 0xD10000
	general = 0x7EC8E3
)

func Handle(eventType string, body []byte) (*discordgo.MessageEmbed, error) {
	switch eventType {
	case "repo:push":
		return handlePush(body)
	case "repo:commit_status_created":
		return commitStatusCreated(body)
	case "repo:commit_status_updated":
		return commitStatusUpdated(body)
	case "pullrequest:created":
		return pullRequestCreated(body)
	case "pullrequest:updated":
		return pullRequestUpdated(body)
	case "pullrequest:approved":
		return pullRequestApproved(body)
	case "pullrequest:unapproved":
		return pullRequestUnapproved(body)
	case "pullrequest:fulfilled":
		return pullRequestFulfilled(body)
	case "pullrequest:rejected":
		return pullRequestRejected(body)
	case "pullrequest:comment_created":
		return pullRequestCommentCreated(body)
	case "pullrequest:comment_updated":
		return pullRequestCommentUpdated(body)
	case "pullrequest:comment_deleted":
		return pullRequestCommentDeleted(body)
	}
	return nil, nil
}

func handlePush(body []byte) (*discordgo.MessageEmbed, error) {
	var push RepoPushEvent
	err := json.Unmarshal(body, &push)
	if err != nil {
		return nil, err
	}
	numOfCommits := 0
	resourceName := "unknown"
	resourceType := "unknown"
	if len(push.Push.Changes) > 0 {
		numOfCommits = len(push.Push.Changes[0].Commits)
		resourceName = push.Push.Changes[0].New.Name
		resourceType = push.Push.Changes[0].New.Type
	}

	return embed.NewEmbed().
		SetTitle(push.Repository.Name+" - Push happened").
		AddField("Number of commits", fmt.Sprintf("%d", numOfCommits)).
		AddField("Resource name", resourceName).
		AddField("Resource type", resourceType).
		SetColor(success).
		SetDescription(push.Actor.DisplayName + " pushed").
		MessageEmbed, nil
}

func commitStatusCreated(body []byte) (*discordgo.MessageEmbed, error) {
	var push RepoPushEvent
	err := json.Unmarshal(body, &push)
	if err != nil {
		return nil, err
	}
	return embed.NewEmbed().SetTitle(push.Repository.Name + " - Push").SetDescription(push.Actor.DisplayName + "pushed").MessageEmbed, nil
}

func commitStatusUpdated(body []byte) (*discordgo.MessageEmbed, error) {
	var push RepoPushEvent
	err := json.Unmarshal(body, &push)
	if err != nil {
		return nil, err
	}
	return embed.NewEmbed().SetTitle(push.Repository.Name + " - Push").SetDescription(push.Actor.DisplayName + "pushed").MessageEmbed, nil
}

func pullRequestCreated(body []byte) (*discordgo.MessageEmbed, error) {
	var created PullRequestCreatedEvent
	err := json.Unmarshal(body, &created)
	if err != nil {
		return nil, err
	}
	return embed.NewEmbed().
		SetTitle(created.Repository.Name + " - Pull request created").
		SetDescription(created.PullRequest.Source.Branch.Name+" -> "+created.PullRequest.Destination.Branch.Name).
		SetURL(created.PullRequest.Links.HTML.Href).
		AddField("Author", created.Actor.DisplayName).
		AddField("Name", created.PullRequest.Title).
		AddField("Reviewers", "").
		SetColor(general).
		MessageEmbed, nil
}

func pullRequestUpdated(body []byte) (*discordgo.MessageEmbed, error) {
	var push RepoPushEvent
	err := json.Unmarshal(body, &push)
	if err != nil {
		return nil, err
	}
	return embed.NewEmbed().SetTitle(push.Repository.Name + " - Push").SetDescription(push.Actor.DisplayName + "pushed").MessageEmbed, nil
}

func pullRequestApproved(body []byte) (*discordgo.MessageEmbed, error) {
	var push RepoPushEvent
	err := json.Unmarshal(body, &push)
	if err != nil {
		return nil, err
	}
	return embed.NewEmbed().SetTitle(push.Repository.Name + " - Push").SetDescription(push.Actor.DisplayName + "pushed").MessageEmbed, nil
}

func pullRequestUnapproved(body []byte) (*discordgo.MessageEmbed, error) {
	var push RepoPushEvent
	err := json.Unmarshal(body, &push)
	if err != nil {
		return nil, err
	}
	return embed.NewEmbed().SetTitle(push.Repository.Name + " - Push").SetDescription(push.Actor.DisplayName + "pushed").MessageEmbed, nil
}

func pullRequestFulfilled(body []byte) (*discordgo.MessageEmbed, error) {
	var push RepoPushEvent
	err := json.Unmarshal(body, &push)
	if err != nil {
		return nil, err
	}
	return embed.NewEmbed().SetTitle(push.Repository.Name + " - Push").SetDescription(push.Actor.DisplayName + "pushed").MessageEmbed, nil
}

func pullRequestRejected(body []byte) (*discordgo.MessageEmbed, error) {
	var push RepoPushEvent
	err := json.Unmarshal(body, &push)
	if err != nil {
		return nil, err
	}
	return embed.NewEmbed().SetTitle(push.Repository.Name + " - Push").SetDescription(push.Actor.DisplayName + "pushed").MessageEmbed, nil
}

func pullRequestCommentCreated(body []byte) (*discordgo.MessageEmbed, error) {
	var push RepoPushEvent
	err := json.Unmarshal(body, &push)
	if err != nil {
		return nil, err
	}
	return embed.NewEmbed().SetTitle(push.Repository.Name + " - Push").SetDescription(push.Actor.DisplayName + "pushed").MessageEmbed, nil
}

func pullRequestCommentUpdated(body []byte) (*discordgo.MessageEmbed, error) {
	var push RepoPushEvent
	err := json.Unmarshal(body, &push)
	if err != nil {
		return nil, err
	}
	return embed.NewEmbed().SetTitle(push.Repository.Name + " - Push").SetDescription(push.Actor.DisplayName + "pushed").MessageEmbed, nil
}

func pullRequestCommentDeleted(body []byte) (*discordgo.MessageEmbed, error) {
	var push RepoPushEvent
	err := json.Unmarshal(body, &push)
	if err != nil {
		return nil, err
	}
	return embed.NewEmbed().SetTitle(push.Repository.Name + " - Push").SetDescription(push.Actor.DisplayName + "pushed").MessageEmbed, nil
}
