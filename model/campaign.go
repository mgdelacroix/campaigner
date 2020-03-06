package model

import (
	"fmt"
	"io"
)

// ToDo: add key-value extra params as a map to allow for customfield_whatever = team
type Campaign struct {
	Url       string    `json:"url"`
	Project   string    `json:"project"`
	Epic      string    `json:"epic"`
	IssueType string    `json:"issue_type"`
	Summary   string    `json:"summary"`
	Template  string    `json:"template"`
	Tickets   []*Ticket `json:"tickets,omitempty"`
}

func (c *Campaign) NextUnpublishedTicket() *Ticket {
	for _, ticket := range c.Tickets {
		if ticket.JiraLink == "" {
			return ticket
		}
	}
	return nil
}

func (c *Campaign) PrintStatus(w io.Writer) {
	fmt.Fprintf(w, "Url: %s\n", c.Url)
	fmt.Fprintf(w, "Project: %s\n", c.Project)
	fmt.Fprintf(w, "Epic: %s\n", c.Epic)
	fmt.Fprintf(w, "Issue Type: %s\n", c.IssueType)
	fmt.Fprintf(w, "Summary: %s\n", c.Summary)
	fmt.Fprintf(w, "Template: %s\n", c.Template)

	for _, ticket := range c.Tickets {
		ticket.PrintStatus(w)
	}
}
