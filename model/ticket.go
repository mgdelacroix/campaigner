package model

import (
	"fmt"
	"io"
)

type Ticket struct {
	GithubLink   int64                  `json:"github_link,omitempty"`
	GithubStatus string                 `json:"github_status,omitempty"`
	JiraLink     string                 `json:"jira_link,omitempty"`
	JiraStatus   string                 `json:"jira_status,omitempty"`
	Summary      string                 `json:"summary,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
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

func (t *Ticket) PrintStatus(w io.Writer) {
	if t.Summary != "" {
		fmt.Fprintf(w, "[%s] %s\n", t.JiraLink, t.Summary)
	}
}
