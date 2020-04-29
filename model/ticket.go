package model

import (
	"fmt"
)

type Ticket struct {
	GithubLink     int                    `json:"github_link,omitempty"`
	GithubStatus   string                 `json:"github_status,omitempty"`
	GithubAssignee string                 `json:"github_assignee,omitempty"`
	JiraLink       string                 `json:"jira_link,omitempty"`
	JiraStatus     string                 `json:"jira_status,omitempty"`
	Summary        string                 `json:"summary,omitempty"`
	Description    string                 `json:"description,omitempty"`
	Data           map[string]interface{} `json:"data,omitempty"`
}

func (t *Ticket) IsPublishedJira() bool {
	return t.JiraLink != ""
}

func (t *Ticket) IsPublishedGithub() bool {
	return t.JiraLink != "" && t.GithubLink != 0
}

func (t *Ticket) PrintStatus() {
	if t.Summary != "" {
		fmt.Printf("[%s] %s\n", t.JiraLink, t.Summary)
	}
}
