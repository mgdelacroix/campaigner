package model

import (
	"bytes"
	"fmt"
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
		if !ticket.IsPublishedJira() {
			return ticket
		}
	}
	return nil
}

func (c *Campaign) NextGithubUnpublishedTicket() *Ticket {
	for _, ticket := range c.Tickets {
		if ticket.IsPublishedJira() && !ticket.IsPublishedGithub() {
			return ticket
		}
	}
	return nil
}

func (c *Campaign) PrintStatus() {
	fmt.Printf("JIRA URL: %s\n", c.Jira.Url)
	fmt.Printf("JIRA Project: %s\n", c.Jira.Project)
	fmt.Printf("JIRA Epic: %s\n", c.Jira.Epic)
	fmt.Printf("JIRA Issue Type: %s\n", c.Jira.IssueType)
	fmt.Printf("GitHub Repo: %s\n", c.Github.Repo)
	fmt.Printf("GitHub Labels: %s\n", c.Github.Labels)
	fmt.Printf("Summary: %s\n", c.Summary)
	fmt.Printf("Template: %s\n", c.Template)
	fmt.Println("")

	for _, ticket := range c.Tickets {
		ticket.PrintStatus()
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
