package model

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/fatih/color"
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
	Jira           ConfigJira   `json:"jira"`
	Github         ConfigGithub `json:"github"`
	Summary        string       `json:"summary"`
	IssueTemplate  string       `json:"issue_template"`
	FooterTemplate string       `json:"footer_template"`
	Tickets        []*Ticket    `json:"tickets,omitempty"`
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
	totalTickets := len(c.Tickets)
	var totalPublishedJira, totalPublishedGithub int
	for _, t := range c.Tickets {
		if t.IsPublishedJira() {
			totalPublishedJira++
			if t.IsPublishedGithub() {
				totalPublishedGithub++
			}
		}
	}

	fmt.Printf("Current campaign for %s with summary\n%s\n\n", color.GreenString(c.Github.Repo), color.CyanString(c.Summary))
	fmt.Printf("\t%d\ttotal tickets\n", totalTickets)
	fmt.Printf("\t%d/%d\tpublished in Jira\n", totalPublishedJira, totalTickets)
	fmt.Printf("\t%d/%d\tpublished in Github\n\n", totalPublishedGithub, totalPublishedJira)
}

func (c *Campaign) AddTickets(tickets []*Ticket, fileOnly bool) {
	c.Tickets = append(c.Tickets, tickets...)
	c.RemoveDuplicateTickets(fileOnly)
}

func (c *Campaign) RemoveDuplicateTickets(fileOnly bool) {
	datalessTickets := []*Ticket{}
	ticketMap := map[string]*Ticket{}
	for _, t := range c.Tickets {
		filename, _ := t.Data["filename"].(string)
		lineNo, _ := t.Data["lineNo"].(int)

		if filename == "" {
			datalessTickets = append(datalessTickets, t)
			continue
		}

		if fileOnly {
			ticketMap[filename] = t
		} else {
			ticketMap[fmt.Sprintf("%s:%d", filename, lineNo)] = t
		}
	}

	cleanTickets := []*Ticket{}
	// dataless tickets are added first as they come from already
	// existing tickets in Jira
	cleanTickets = append(cleanTickets, datalessTickets...)
	for _, t := range ticketMap {
		cleanTickets = append(cleanTickets, t)
	}

	c.Tickets = cleanTickets
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

	descriptionTemplate, err := template.ParseFiles(c.IssueTemplate)
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

func (c *Campaign) GetJiraUrl(ticket *Ticket) string {
	return fmt.Sprintf("%s/browse/%s", c.Jira.Url, ticket.JiraLink)
}

func (c *Campaign) GetGithubUrl(ticket *Ticket) string {
	return fmt.Sprintf("https://github.com/%s/issues/%d", c.Github.Repo, ticket.GithubLink)
}
