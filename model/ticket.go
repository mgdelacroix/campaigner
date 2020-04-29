package model

import (
	"fmt"
	"strings"
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

func RemoveDuplicateTickets(tickets []*Ticket, fileOnly bool) []*Ticket {
	ticketMap := map[string]*Ticket{}
	for _, t := range tickets {
		filename, _ := t.Data["filename"].(string)
		lineNo, _ := t.Data["lineNo"].(int)
		if fileOnly {
			ticketMap[filename] = t
		} else {
			ticketMap[fmt.Sprintf("%s:%d", filename, lineNo)] = t
		}
	}

	cleanTickets := []*Ticket{}
	for _, t := range ticketMap {
		cleanTickets = append(cleanTickets, t)
	}

	return cleanTickets
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

func (t *Ticket) JiraIssue() string {
	parts := strings.Split(t.JiraLink, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-1]
}
