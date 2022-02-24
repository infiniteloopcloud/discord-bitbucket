package bitbucket

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"github.com/infiniteloopcloud/discord-bitbucket/env"
)

const (
	success   = 0x90EE90
	failure   = 0xD10000
	prCreated = 0x89CFF0
	prUpdated = 0x0047AB
	gray      = 0x979797
)

func Handle(eventType string, body []byte) (string, *discordgo.MessageEmbed, error) {
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
	return "", nil, nil
}

func handlePush(body []byte) (string, *discordgo.MessageEmbed, error) {
	if os.Getenv(env.SkipRepoPushMessages) == "true" {
		return "", nil, nil
	}

	var push RepoPushEvent
	err := json.Unmarshal(body, &push)
	if err != nil {
		return "", nil, err
	}
	numOfCommits := 0
	resourceName := "unknown"
	resourceType := "unknown"
	if len(push.Push.Changes) > 0 {
		numOfCommits = len(push.Push.Changes[0].Commits)
		resourceName = push.Push.Changes[0].New.Name
		resourceType = push.Push.Changes[0].New.Type
	}

	message := embed.NewEmbed().
		SetTitle("Push happened").
		AddField("Number of commits", fmt.Sprintf("%d", numOfCommits)).
		AddField("Resource name", resourceName).
		AddField("Resource type", resourceType).
		SetColor(success)
	if push.Actor.DisplayName != "" {
		message = message.SetDescription(push.Actor.DisplayName + " pushed")
	}

	return push.Repository.Name, message.MessageEmbed, nil
}

func commitStatusCreated(body []byte) (string, *discordgo.MessageEmbed, error) {
	var push RepoPushEvent
	err := json.Unmarshal(body, &push)
	if err != nil {
		return "", nil, err
	}
	return "", embed.NewEmbed().SetTitle(push.Repository.Name + " - Push").SetDescription(push.Actor.DisplayName + "pushed").MessageEmbed, nil
}

func commitStatusUpdated(body []byte) (string, *discordgo.MessageEmbed, error) {
	var event RepoCommitStatusUpdatedEvent
	err := json.Unmarshal(body, &event)
	if err != nil {
		return "", nil, err
	}
	color := gray
	if event.CommitStatus.State == "FAILED" {
		color = failure
	} else if event.CommitStatus.State == "SUCCESSFUL" {
		color = success
	}

	if event.CommitStatus.Name == "" {
		return "", nil, nil
	}

	message := embed.NewEmbed().SetTitle(event.CommitStatus.Name).SetColor(color).SetDescription("Pipeline trigger")

	if event.CommitStatus.State != "" {
		message = message.AddField("Status", event.CommitStatus.State)
	}
	if event.CommitStatus.Commit.Author.User.DisplayName != "" {
		message = message.AddField("Triggered by", event.CommitStatus.Commit.Author.User.DisplayName)
	}
	if event.CommitStatus.URL != "" {
		message = message.SetURL(event.CommitStatus.URL)
	}

	return event.Repository.Name, message.MessageEmbed, nil
}

func pullRequestCreated(body []byte) (string, *discordgo.MessageEmbed, error) {
	var created PullRequestCreatedEvent
	err := json.Unmarshal(body, &created)
	if err != nil {
		return "", nil, err
	}
	reviewers := "none"
	reviewerList := []string{}
	for _, reviewer := range created.PullRequest.Reviewers {
		reviewerList = append(reviewerList, reviewer.DisplayName)
	}
	if len(reviewerList) > 0 {
		reviewers = strings.Join(reviewerList, ", ")
	}

	if created.Actor.DisplayName == "" || created.PullRequest.Title == "" {
		return "", nil, nil
	}

	message := embed.NewEmbed().
		SetTitle(created.Actor.DisplayName+"created a new pull request: "+created.PullRequest.Title).SetColor(prCreated).AddField("Reviewers", reviewers)

	if created.PullRequest.Source.Branch.Name != "" && created.PullRequest.Destination.Branch.Name != "" {
		message = message.SetDescription("`" + created.PullRequest.Source.Branch.Name + "` > `" + created.PullRequest.Destination.Branch.Name + "`")
	}
	if created.PullRequest.Links.HTML.Href != "" {
		message = message.SetURL(created.PullRequest.Links.HTML.Href)
	}
	if created.PullRequest.State != "" {
		message = message.AddField("Status", created.PullRequest.State)
	}
	if created.PullRequest.Description != "" {
		if len(created.PullRequest.Description) > 200 {
			desc := created.PullRequest.Description[0:199] + "..."
			message = message.AddField("PR Description", desc)
		} else {
			message = message.AddField("PR Description", created.PullRequest.Description)
		}
	}

	return created.Repository.Name, message.MessageEmbed, nil
}

func pullRequestUpdated(body []byte) (string, *discordgo.MessageEmbed, error) {
	var updated PullRequestUpdatedEvent
	err := json.Unmarshal(body, &updated)
	if err != nil {
		return "", nil, err
	}
	reviewers := "none"
	reviewerList := []string{}
	for _, reviewer := range updated.PullRequest.Reviewers {
		reviewerList = append(reviewerList, reviewer.DisplayName)
	}
	if len(reviewerList) > 0 {
		reviewers = strings.Join(reviewerList, ", ")
	}

	if updated.Actor.DisplayName == "" || updated.PullRequest.Title == "" {
		return "", nil, nil
	}

	message := embed.NewEmbed().
		SetTitle(updated.Actor.DisplayName+"updated the pull request: "+updated.PullRequest.Title).SetColor(prCreated).AddField("Reviewers", reviewers)

	if updated.PullRequest.Source.Branch.Name != "" && updated.PullRequest.Destination.Branch.Name != "" {
		message = message.SetDescription("`" + updated.PullRequest.Source.Branch.Name + "` > `" + updated.PullRequest.Destination.Branch.Name + "`")
	}
	if updated.PullRequest.Links.HTML.Href != "" {
		message = message.SetURL(updated.PullRequest.Links.HTML.Href)
	}
	if updated.PullRequest.State != "" {
		message = message.AddField("Status", updated.PullRequest.State)
	}
	if updated.PullRequest.Description != "" {
		if len(updated.PullRequest.Description) > 200 {
			desc := updated.PullRequest.Description[0:199] + "..."
			message = message.AddField("PR Description", desc)
		} else {
			message = message.AddField("PR Description", updated.PullRequest.Description)
		}
	}

	return updated.Repository.Name, message.MessageEmbed, nil
}

func pullRequestApproved(body []byte) (string, *discordgo.MessageEmbed, error) {
	var approved PullRequestApprovedEvent
	err := json.Unmarshal(body, &approved)
	if err != nil {
		return "", nil, err
	}

	if approved.Actor.DisplayName == "" || approved.PullRequest.Title == "" {
		return "", nil, nil
	}

	message := embed.NewEmbed().SetTitle(approved.Approval.User.DisplayName + " approved pull request: " + approved.PullRequest.Title).SetColor(success)

	if approved.PullRequest.Source.Branch.Name != "" && approved.PullRequest.Destination.Branch.Name != "" {
		message = message.SetDescription("`" + approved.PullRequest.Source.Branch.Name + "` > `" + approved.PullRequest.Destination.Branch.Name + "`")
	}
	if approved.PullRequest.Links.HTML.Href != "" {
		message = message.SetURL(approved.PullRequest.Links.HTML.Href)
	}
	if approved.Actor.DisplayName != "" {
		message = message.AddField("Created by", approved.Actor.DisplayName)
	}

	return approved.Repository.Name, message.MessageEmbed, nil
}

func pullRequestUnapproved(body []byte) (string, *discordgo.MessageEmbed, error) {
	var unapproved PullRequestApprovedEvent
	err := json.Unmarshal(body, &unapproved)
	if err != nil {
		return "", nil, err
	}

	if unapproved.Actor.DisplayName == "" || unapproved.PullRequest.Title == "" {
		return "", nil, nil
	}

	message := embed.NewEmbed().SetTitle(unapproved.Approval.User.DisplayName + " unapproved pull request: " + unapproved.PullRequest.Title).SetColor(success)

	if unapproved.PullRequest.Source.Branch.Name != "" && unapproved.PullRequest.Destination.Branch.Name != "" {
		message = message.SetDescription("`" + unapproved.PullRequest.Source.Branch.Name + "` > `" + unapproved.PullRequest.Destination.Branch.Name + "`")
	}
	if unapproved.PullRequest.Links.HTML.Href != "" {
		message = message.SetURL(unapproved.PullRequest.Links.HTML.Href)
	}
	if unapproved.Actor.DisplayName != "" {
		message = message.AddField("Created by", unapproved.Actor.DisplayName)
	}

	return unapproved.Repository.Name, message.MessageEmbed, nil
}

func pullRequestFulfilled(body []byte) (string, *discordgo.MessageEmbed, error) {
	var merged PullRequestMergedEvent
	err := json.Unmarshal(body, &merged)
	if err != nil {
		return "", nil, err
	}
	reviewers := "none"
	reviewerList := []string{}
	for _, reviewer := range merged.PullRequest.Reviewers {
		reviewerList = append(reviewerList, reviewer.DisplayName)
	}
	if len(reviewerList) > 0 {
		reviewers = strings.Join(reviewerList, ", ")
	}

	if merged.PullRequest.ClosedBy.DisplayName == "" || merged.PullRequest.Title == "" {
		return "", nil, nil
	}

	message := embed.NewEmbed().SetTitle(merged.PullRequest.ClosedBy.DisplayName+" merged pull request: "+merged.PullRequest.Title).SetColor(success).AddField("Approved by", reviewers)

	if merged.PullRequest.Source.Branch.Name != "" && merged.PullRequest.Destination.Branch.Name != "" {
		message = message.SetDescription("`" + merged.PullRequest.Source.Branch.Name + "` > `" + merged.PullRequest.Destination.Branch.Name + "`")
	}
	if merged.PullRequest.Links.HTML.Href != "" {
		message = message.SetURL(merged.PullRequest.Links.HTML.Href)
	}
	if merged.Actor.DisplayName != "" {
		message = message.AddField("Created by", merged.Actor.DisplayName)
	}
	if merged.PullRequest.State != "" {
		message = message.AddField("Status", merged.PullRequest.State)
	}

	return merged.Repository.Name, message.MessageEmbed, nil
}

func pullRequestRejected(body []byte) (string, *discordgo.MessageEmbed, error) {
	var rejected PullRequestMergedEvent
	err := json.Unmarshal(body, &rejected)
	if err != nil {
		return "", nil, err
	}

	if rejected.PullRequest.ClosedBy.DisplayName == "" || rejected.PullRequest.Title == "" {
		return "", nil, nil
	}

	message := embed.NewEmbed().SetTitle(rejected.PullRequest.ClosedBy.DisplayName + " declined pull request: " + rejected.PullRequest.Title).SetColor(failure)

	if rejected.PullRequest.Source.Branch.Name != "" && rejected.PullRequest.Destination.Branch.Name != "" {
		message = message.SetDescription("`" + rejected.PullRequest.Source.Branch.Name + "` > `" + rejected.PullRequest.Destination.Branch.Name + "`")
	}
	if rejected.PullRequest.Links.HTML.Href != "" {
		message = message.SetURL(rejected.PullRequest.Links.HTML.Href)
	}
	if rejected.Actor.DisplayName != "" {
		message = message.AddField("Created by", rejected.Actor.DisplayName)
	}

	return rejected.Repository.Name, message.MessageEmbed, nil
}

func pullRequestCommentCreated(body []byte) (string, *discordgo.MessageEmbed, error) {
	var commentCreated PullRequestCommentCreatedEvent
	err := json.Unmarshal(body, &commentCreated)
	if err != nil {
		return "", nil, err
	}

	comment := "no comment"
	if commentCreated.Comment.Content.Raw != "" {
		comment = commentCreated.Comment.Content.Raw
		if len(comment) > 105 {
			comment = commentCreated.Comment.Content.Raw[0:100] + "..."
		}
	}

	if commentCreated.Comment.User.DisplayName == "" || commentCreated.PullRequest.Title == "" {
		return "", nil, nil
	}

	message := embed.NewEmbed().
		SetTitle(commentCreated.Comment.User.DisplayName+" commented pull request: "+commentCreated.PullRequest.Title).
		AddField("Comment", comment).
		SetColor(prCreated)

	if commentCreated.PullRequest.Source.Branch.Name != "" && commentCreated.PullRequest.Destination.Branch.Name != "" {
		message = message.SetDescription("`" + commentCreated.PullRequest.Source.Branch.Name + "` > `" + commentCreated.PullRequest.Destination.Branch.Name + "`")
	}
	if commentCreated.Comment.Links.HTML.Href != "" {
		message = message.SetURL(commentCreated.Comment.Links.HTML.Href)
	}

	return commentCreated.Repository.Name, message.MessageEmbed, nil
}

func pullRequestCommentUpdated(_ []byte) (string, *discordgo.MessageEmbed, error) {
	return "", nil, errors.New("not supported: pullRequestCommentUpdated")
}

func pullRequestCommentDeleted(body []byte) (string, *discordgo.MessageEmbed, error) {
	var commentDeleted PullRequestCommentCreatedEvent
	err := json.Unmarshal(body, &commentDeleted)
	if err != nil {
		return "", nil, err
	}

	if commentDeleted.Comment.User.DisplayName == "" || commentDeleted.PullRequest.Title == "" {
		return "", nil, nil
	}

	message := embed.NewEmbed().
		SetTitle(commentDeleted.Comment.User.DisplayName + " comment deleted on pull request: " + commentDeleted.PullRequest.Title).
		SetColor(failure)

	if commentDeleted.PullRequest.Source.Branch.Name != "" && commentDeleted.PullRequest.Destination.Branch.Name != "" {
		message = message.SetDescription("`" + commentDeleted.PullRequest.Source.Branch.Name + "` > `" + commentDeleted.PullRequest.Destination.Branch.Name + "`")
	}
	if commentDeleted.Comment.Links.HTML.Href != "" {
		message = message.SetURL(commentDeleted.Comment.Links.HTML.Href)
	}

	return commentDeleted.Repository.Name, message.MessageEmbed, nil
}
