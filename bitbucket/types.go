package bitbucket

import "time"

// RepoPushEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-Push
type RepoPushEvent struct {
	Actor      Actor      `json:"actor"`
	Repository Repository `json:"repository"`
	Push       struct {
		Changes []struct {
			Forced    bool     `json:"forced"`
			Old       OldOrNew `json:"old"`
			New       OldOrNew `json:"new"`
			Closed    bool     `json:"closed"`
			Created   bool     `json:"created"`
			Truncated bool     `json:"truncated"`
			Links     `json:"links"`
			Commits   []Commit `json:"commits"`
		} `json:"changes"`
	} `json:"push"`
}

// Links is a common struct used in several types. Refer to the event documentation
// to find out which link types are populated in which events.
type Links struct {
	Avatar struct {
		Href string `json:"href"`
	} `json:"avatar"`
	HTML struct {
		Href string `json:"href"`
	} `json:"html"`
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
	Commits struct {
		Href string `json:"href"`
	} `json:"commits"`
	Commit struct {
		Href string `json:"href"`
	} `json:"commit"`
}

// OldOrNew is used in the RepoPushEvent type
type OldOrNew struct {
	Repository struct {
		FullName string `json:"full_name"`
		UUID     string `json:"uuid"`
		Links    Links  `json:"links"`
		Name     string `json:"name"`
		Type     string `json:"type"`
	} `json:"repository"`
	Target struct {
		Date    *time.Time `json:"date"`
		Parents []struct {
			Hash  string `json:"hash"`
			Links Links  `json:"links"`
			Type  string `json:"type"`
		} `json:"parents"`
		Message string `json:"message"`
		Hash    string `json:"hash"`
		Author  Author `json:"author"`
		Links   Links  `json:"links"`
		Type    string `json:"type"`
	} `json:"target"`
	Links Links  `json:"links"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

// Author is a common struct used in several types
type Author struct {
	Raw  string `json:"raw"`
	User struct {
		Username    string `json:"username"`
		Type        string `json:"type"`
		UUID        string `json:"uuid"`
		Links       Links  `json:"links"`
		DisplayName string `json:"display_name"`
	} `json:"user"`
}

// Repository is a common struct used in several types
type Repository struct {
	Scm      string `json:"scm"`
	FullName string `json:"full_name"`
	Type     string `json:"type"`
	Website  string `json:"website"`
	Owner    struct {
		Username    string `json:"username"`
		Type        string `json:"type"`
		UUID        string `json:"uuid"`
		Links       Links  `json:"links"`
		DisplayName string `json:"display_name"`
	} `json:"owner"`
	UUID      string `json:"uuid"`
	Links     Links  `json:"links"`
	Name      string `json:"name"`
	IsPrivate bool   `json:"is_private"`
}

// Commit is a common struct used in several types
type Commit struct {
	Date    time.Time `json:"date"`
	Parents []struct {
		Hash  string `json:"hash"`
		Links Links  `json:"self"`
		Type  string `json:"type"`
	} `json:"parents"`
	Message string `json:"message"`
	Hash    string `json:"hash"`
	Author  Author `json:"author"`
	Links   Links  `json:"links"`
	Type    string `json:"type"`
}

// Actor is a common struct used in several types
type Actor struct {
	Nickname    string `json:"nickname"`
	Type        string `json:"type"`
	UUID        string `json:"uuid"`
	Links       Links  `json:"links"`
	DisplayName string `json:"display_name"`
}

// RepoForkEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-Fork
type RepoForkEvent struct {
	Actor      Actor      `json:"actor"`
	Repository Repository `json:"repository"`
	Fork       Repository `json:"fork"`
}

// Comment https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-entity_comment
type Comment struct {
	ID     int `json:"id"`
	Parent struct {
		ID int `json:"id"`
	} `json:"parent"`
	Content struct {
		Raw    string `json:"raw"`
		HTML   string `json:"html"`
		Markup string `json:"markup"`
	} `json:"content"`
	Inline struct {
		Path string      `json:"path"`
		From interface{} `json:"from"`
		To   int         `json:"to"`
	} `json:"inline"`
	User      User       `json:"user"`
	CreatedOn *time.Time `json:"created_on"`
	UpdatedOn *time.Time `json:"updated_on"`
	Links     Links      `json:"links"`
}

// RepoCommitCommentCreatedEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-CommitCommentCreated
type RepoCommitCommentCreatedEvent struct {
	Actor      Actor      `json:"actor"`
	Comment    Comment    `json:"comment"`
	Repository Repository `json:"repository"`
	Commit     struct {
		Hash string `json:"hash"`
	} `json:"commit"`
}

// A RepoCommitStatusEvent is not a BB event. This is the base for several CommitStatus* events.
type RepoCommitStatusEvent struct {
	Actor        Actor      `json:"actor"`
	Repository   Repository `json:"repository"`
	CommitStatus struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		State       string `json:"state"`
		Key         string `json:"key"`
		URL         string `json:"url"`
		Type        string `json:"type"`
		Commit      struct {
			Hash   string `json:"hash"`
			Author struct {
				User User `json:"user"`
			} `json:"author"`
		} `json:"commit"`
		CreatedOn *time.Time `json:"created_on"`
		UpdatedOn *time.Time `json:"updated_on"`
		Links     Links      `json:"links"`
	} `json:"commit_status"`
}

// RepoCommitStatusCreatedEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-CommitStatusCreated
type RepoCommitStatusCreatedEvent struct {
	RepoCommitStatusEvent
}

// RepoCommitStatusUpdatedEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-CommitStatusUpdated
type RepoCommitStatusUpdatedEvent struct {
	RepoCommitStatusEvent
}

// An IssueEvent is not a BB event. This is the base for several Issue* events.
type IssueEvent struct {
	Actor      Actor      `json:"actor"`
	Issue      Issue      `json:"issue"`
	Repository Repository `json:"repository"`
}

// IssueCreatedEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-Created
type IssueCreatedEvent struct {
	IssueEvent
}

// IssueUpdatedEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-Updated
type IssueUpdatedEvent struct {
	IssueEvent
	Comment Comment `json:"comment"`
	Changes struct {
		Status struct {
			Old string `json:"old"`
			New string `json:"new"`
		} `json:"status"`
	} `json:"changes"`
}

// IssueCommentCreatedEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-CommentCreated
type IssueCommentCreatedEvent struct {
	IssueEvent
	Comment Comment `json:"comment"`
}

// Issue https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-entity_issue
type Issue struct {
	ID        int    `json:"id"`
	Component string `json:"component"`
	Title     string `json:"title"`
	Content   struct {
		Raw    string `json:"raw"`
		HTML   string `json:"html"`
		Markup string `json:"markup"`
	} `json:"content"`
	Priority  string `json:"priority"`
	State     string `json:"state"`
	Type      string `json:"type"`
	Milestone struct {
		Name string `json:"name"`
	} `json:"milestone"`
	Version struct {
		Name string `json:"name"`
	} `json:"version"`
	CreatedOn *time.Time `json:"created_on"`
	UpdatedOn *time.Time `json:"updated_on"`
	Links     Links      `json:"links"`
}

// A PullRequestEvent is not a BB event. This is the base for several PullRequest* events.
type PullRequestEvent struct {
	Actor       Actor       `json:"actor"`
	PullRequest PullRequest `json:"pullrequest"`
	Repository  Repository  `json:"repository"`
}

// PullRequestCreatedEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-Created
type PullRequestCreatedEvent struct {
	PullRequestEvent
}

// PullRequestUpdatedEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-Updated.1
type PullRequestUpdatedEvent struct {
	PullRequestEvent
}

// PullRequestApprovedEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-Approved
type PullRequestApprovedEvent struct {
	PullRequestEvent
	Approval Approval `json:"approval"`
}

// PullRequestApprovalRemovedEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-ApprovalRemoved
type PullRequestApprovalRemovedEvent struct {
	PullRequestEvent
	Approval Approval `json:"approval"`
}

// PullRequestMergedEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-Merged
type PullRequestMergedEvent struct {
	PullRequestEvent
}

// PullRequestDeclinedEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-Declined
type PullRequestDeclinedEvent struct {
	PullRequestEvent
}

// A PullRequestCommentEvent doesn't exist. It is used as the base for several real events.
type PullRequestCommentEvent struct {
	PullRequestEvent
	Comment Comment `json:"comment"`
}

// PullRequestCommentCreatedEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-CommentCreated.1
type PullRequestCommentCreatedEvent struct {
	PullRequestCommentEvent
}

// PullRequestCommentUpdatedEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-CommentUpdated
type PullRequestCommentUpdatedEvent struct {
	PullRequestCommentEvent
}

// PullRequestCommentDeletedEvent https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-CommentDeleted
type PullRequestCommentDeletedEvent struct {
	PullRequestCommentEvent
}

// An Approval is used in pull requests
type Approval struct {
	Date *time.Time `json:"date"`
	User User       `json:"user"`
}

// PullRequest https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-entity_pullrequest
type PullRequest struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"`
	Author      User   `json:"author"`
	Source      struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
		Commit struct {
			Hash string `json:"hash"`
		} `json:"commit"`
		Repository Repository `json:"repository"`
	} `json:"source"`
	Destination struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
		Commit struct {
			Hash string `json:"hash"`
		} `json:"commit"`
		Repository Repository `json:"repository"`
	} `json:"destination"`
	MergeCommit struct {
		Hash string `json:"hash"`
	} `json:"merge_commit"`
	Participants      []Participant `json:"participants"`
	Reviewers         []User        `json:"reviewers"`
	CloseSourceBranch bool          `json:"close_source_branch"`
	ClosedBy          User          `json:"closed_by"`
	Reason            string        `json:"reason"`
	CreatedOn         *time.Time    `json:"created_on"`
	UpdatedOn         *time.Time    `json:"updated_on"`
	Links             Links         `json:"links"`
}

// Participant is the actual structure returned in PullRequest events
// Note: this doesn't match the docs!?
type Participant struct {
	Role     string `json:"role"`
	Type     string `json:"type"`
	Approved bool   `json:"approved"`
	User     User   `json:"user"`
}

// User https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html#EventPayloads-entity_userUser
type User struct {
	Type        string `json:"type"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	UUID        string `json:"uuid"`
	Links       Links  `json:"links"`
}
