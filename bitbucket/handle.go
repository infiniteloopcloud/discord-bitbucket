package bitbucket

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

const (
	success   = 0x90EE90
	failure   = 0xD10000
	prCreated = 0x89CFF0
	prUpdated = 0x0047AB
	gray      = 0x979797
)

func Handle(eventType string, body []byte) (*discordgo.MessageEmbed, error) {
	switch eventType {
	case "repo:push":
		return handlePush(body)
	//case "repo:commit_status_created":
	//	return commitStatusCreated(body)
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
	var event RepoCommitStatusUpdatedEvent
	err := json.Unmarshal(body, &event)
	if err != nil {
		return nil, err
	}
	color := gray
	if event.CommitStatus.State == "FAILED" {
		color = failure
	} else if event.CommitStatus.State == "SUCCESSFUL" {
		color = success
	}

	return embed.NewEmbed().
		SetTitle(event.CommitStatus.Name).
		AddField("Status", event.CommitStatus.State).
		AddField("Triggered by", event.CommitStatus.Commit.Author.User.DisplayName).
		SetColor(color).
		SetDescription("Pipeline trigger").
		SetURL(event.CommitStatus.URL).
		MessageEmbed, nil
}

func pullRequestCreated(body []byte) (*discordgo.MessageEmbed, error) {
	var created PullRequestCreatedEvent
	err := json.Unmarshal(body, &created)
	if err != nil {
		return nil, err
	}
	reviewers := "none"
	reviewerList := []string{}
	for _, reviewer := range created.PullRequest.Reviewers {
		reviewerList = append(reviewerList, reviewer.DisplayName)
	}
	if len(reviewerList) > 0 {
		reviewers = strings.Join(reviewerList, ", ")
	}

	return embed.NewEmbed().
		SetTitle("Pull request created: "+created.PullRequest.Title).
		SetDescription("`"+created.PullRequest.Source.Branch.Name+"` > `"+created.PullRequest.Destination.Branch.Name+"`").
		SetURL(created.PullRequest.Links.HTML.Href).
		AddField("Created by", created.Actor.DisplayName).
		AddField("Reviewers", reviewers).
		AddField("Status", created.PullRequest.State).
		AddField("PR Description", created.PullRequest.Description).
		SetColor(prCreated).
		MessageEmbed, nil
}

func pullRequestUpdated(body []byte) (*discordgo.MessageEmbed, error) {
	var updated PullRequestUpdatedEvent
	err := json.Unmarshal(body, &updated)
	if err != nil {
		return nil, err
	}
	reviewers := "none"
	reviewerList := []string{}
	for _, reviewer := range updated.PullRequest.Reviewers {
		reviewerList = append(reviewerList, reviewer.DisplayName)
	}
	if len(reviewerList) > 0 {
		reviewers = strings.Join(reviewerList, ", ")
	}

	return embed.NewEmbed().
		SetTitle("Pull request updated: "+updated.PullRequest.Title).
		SetDescription("`"+updated.PullRequest.Source.Branch.Name+"` > `"+updated.PullRequest.Destination.Branch.Name+"`").
		SetURL(updated.PullRequest.Links.HTML.Href).
		AddField("Created by", updated.Actor.DisplayName).
		AddField("Reviewers", reviewers).
		AddField("Status", updated.PullRequest.State).
		AddField("PR Description", updated.PullRequest.Description).
		SetColor(prUpdated).
		MessageEmbed, nil
}

func pullRequestApproved(body []byte) (*discordgo.MessageEmbed, error) {
	var approved PullRequestApprovedEvent
	err := json.Unmarshal(body, &approved)
	if err != nil {
		return nil, err
	}
	return embed.NewEmbed().
		SetTitle(approved.Approval.User.DisplayName+" approved pull request: "+approved.PullRequest.Title).
		SetDescription("`"+approved.PullRequest.Source.Branch.Name+"` > `"+approved.PullRequest.Destination.Branch.Name+"`").
		SetURL(approved.PullRequest.Links.HTML.Href).
		AddField("Created by", approved.Actor.DisplayName).
		SetColor(success).
		MessageEmbed, nil
}

func pullRequestUnapproved(body []byte) (*discordgo.MessageEmbed, error) {
	var unapproved PullRequestApprovedEvent
	err := json.Unmarshal(body, &unapproved)
	if err != nil {
		return nil, err
	}
	return embed.NewEmbed().
		SetTitle(unapproved.Approval.User.DisplayName+" unapproved pull request: "+unapproved.PullRequest.Title).
		SetDescription("`"+unapproved.PullRequest.Source.Branch.Name+"` > `"+unapproved.PullRequest.Destination.Branch.Name+"`").
		SetURL(unapproved.PullRequest.Links.HTML.Href).
		AddField("Created by", unapproved.Actor.DisplayName).
		SetColor(failure).
		MessageEmbed, nil
}

func pullRequestFulfilled(body []byte) (*discordgo.MessageEmbed, error) {
	var merged PullRequestMergedEvent
	err := json.Unmarshal(body, &merged)
	if err != nil {
		return nil, err
	}
	reviewers := "none"
	reviewerList := []string{}
	for _, reviewer := range merged.PullRequest.Reviewers {
		reviewerList = append(reviewerList, reviewer.DisplayName)
	}
	if len(reviewerList) > 0 {
		reviewers = strings.Join(reviewerList, ", ")
	}

	return embed.NewEmbed().
		SetTitle(merged.PullRequest.ClosedBy.DisplayName+" merged pull request: "+merged.PullRequest.Title).
		SetDescription("`"+merged.PullRequest.Source.Branch.Name+"` > `"+merged.PullRequest.Destination.Branch.Name+"`").
		SetURL(merged.PullRequest.Links.HTML.Href).
		AddField("Created by", merged.Actor.DisplayName).
		AddField("Approved by", reviewers).
		AddField("Status", merged.PullRequest.State).
		SetColor(success).
		MessageEmbed, nil
}

func pullRequestRejected(body []byte) (*discordgo.MessageEmbed, error) {
	var rejected PullRequestMergedEvent
	err := json.Unmarshal(body, &rejected)
	if err != nil {
		return nil, err
	}

	return embed.NewEmbed().
		SetTitle(rejected.PullRequest.ClosedBy.DisplayName+" declined pull request: "+rejected.PullRequest.Title).
		SetDescription("`"+rejected.PullRequest.Source.Branch.Name+"` > `"+rejected.PullRequest.Destination.Branch.Name+"`").
		SetURL(rejected.PullRequest.Links.HTML.Href).
		AddField("Created by", rejected.Actor.DisplayName).
		SetColor(failure).
		MessageEmbed, nil
}

func pullRequestCommentCreated(body []byte) (*discordgo.MessageEmbed, error) {
	var commentCreated PullRequestCommentCreatedEvent
	err := json.Unmarshal(body, &commentCreated)
	if err != nil {
		return nil, err
	}

	comment := commentCreated.Comment.Content.Raw
	if len(comment) > 105 {
		comment = commentCreated.Comment.Content.Raw[0:100] + "..."
	}
	return embed.NewEmbed().
		SetTitle(commentCreated.Comment.User.DisplayName+" commented pull request: "+commentCreated.PullRequest.Title).
		SetDescription("`"+commentCreated.PullRequest.Source.Branch.Name+"` > `"+commentCreated.PullRequest.Destination.Branch.Name+"`").
		SetURL(commentCreated.Comment.Links.HTML.Href).
		AddField("Comment", comment).
		SetColor(prCreated).
		MessageEmbed, nil
}

func pullRequestCommentUpdated(_ []byte) (*discordgo.MessageEmbed, error) {
	return nil, errors.New("not supported: pullRequestCommentUpdated")
}

func pullRequestCommentDeleted(body []byte) (*discordgo.MessageEmbed, error) {
	var commentDeleted PullRequestCommentCreatedEvent
	err := json.Unmarshal(body, &commentDeleted)
	if err != nil {
		return nil, err
	}

	comment := commentDeleted.Comment.Content.Raw
	if len(comment) > 105 {
		comment = commentDeleted.Comment.Content.Raw[0:100] + "..."
	}
	return embed.NewEmbed().
		SetTitle(commentDeleted.Comment.User.DisplayName + " comment deleted on pull request: " + commentDeleted.PullRequest.Title).
		SetDescription("`" + commentDeleted.PullRequest.Source.Branch.Name + "` > `" + commentDeleted.PullRequest.Destination.Branch.Name + "`").
		SetURL(commentDeleted.Comment.Links.HTML.Href).
		SetColor(failure).
		MessageEmbed, nil
}
