package model

import (
	"fmt"
	"io"
)

type Ticket struct {
	GithubLink int64                  `json:"githubLink,omitempty"`
	JiraLink   string                 `json:"jiraLink,omitempty"`
	Summary    string                 `json:"summary,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
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
	fmt.Fprintf(w, " [%s] %s\n", t.JiraLink, t.Summary)
}
