package model

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
)

type ConfigJira struct {
	Url       string `json:"url"`
	Username  string `json:"username"`
	Token     string `json:"token"`
	Project   string `json:"project"`
	Epic      string `json:"epic"`
	IssueType string `json:"issue_type"`
}

type ConfigGithub struct {
	Token  string   `json:"token"`
	Repo   string   `json:"repo"`
	Labels []string `json:"labels"`
}

// ToDo: add key-value extra params as a map to allow for customfield_whatever = team
type Campaign struct {
	Jira     ConfigJira   `json:"jira"`
	Github   ConfigGithub `json:"github"`
	Summary  string       `json:"summary"`
	Template string       `json:"template"`
	Tickets  []*Ticket    `json:"tickets,omitempty"`
}

func (c *Campaign) NextJiraUnpublishedTicket() *Ticket {
	for _, ticket := range c.Tickets {
		if ticket.JiraLink == "" {
			return ticket
		}
	}
	return nil
}

func (c *Campaign) NextGithubUnpublishedTicket() *Ticket {
	for _, ticket := range c.Tickets {
		if ticket.JiraLink != "" && ticket.GithubLink == 0 {
			return ticket
		}
	}
	return nil
}

func (c *Campaign) PrintStatus(w io.Writer) {
	fmt.Fprintf(w, "JIRA URL: %s\n", c.Jira.Url)
	fmt.Fprintf(w, "JIRA Project: %s\n", c.Jira.Project)
	fmt.Fprintf(w, "JIRA Epic: %s\n", c.Jira.Epic)
	fmt.Fprintf(w, "JIRA Issue Type: %s\n", c.Jira.IssueType)
	fmt.Fprintf(w, "GitHub Repo: %s\n", c.Github.Repo)
	fmt.Fprintf(w, "GitHub Labels: %s\n", c.Github.Labels)
	fmt.Fprintf(w, "Summary: %s\n", c.Summary)
	fmt.Fprintf(w, "Template: %s\n", c.Template)
	fmt.Fprintln(w, "")

	for _, ticket := range c.Tickets {
		ticket.PrintStatus(w)
	}
}

func (c *Campaign) FillTicket(t *Ticket) error {
	summaryTmpl, err := template.New("").Parse(c.Summary)
	if err != nil {
		return err
	}

	var summaryBytes bytes.Buffer
	if err := summaryTmpl.Execute(&summaryBytes, t.Data); err != nil {
		return err
	}
	t.Summary = summaryBytes.String()

	descriptionTemplate, err := template.ParseFiles(c.Template)
	if err != nil {
		return err
	}

	var descriptionBytes bytes.Buffer
	if err := descriptionTemplate.Execute(&descriptionBytes, t.Data); err != nil {
		return err
	}
	t.Description = descriptionBytes.String()
	return nil
}

func (c *Campaign) RepoComponents() (string, string) {
	parts := strings.Split(c.Github.Repo, "/")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", ""
}
