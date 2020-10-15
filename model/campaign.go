package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
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
	var totalPublishedJira, totalPublishedGithub, totalAssigned, totalClosed int
	for _, t := range c.Tickets {
		if t.IsPublishedJira() {
			totalPublishedJira++
			if t.IsPublishedGithub() {
				totalPublishedGithub++
				if t.IsAssigned() {
					totalAssigned++
					if t.IsClosed() {
						totalClosed++
					}
				}
			}
		}
	}

	fmt.Printf("Current campaign for %s with summary\n%s\n\n", color.GreenString(c.Github.Repo), color.CyanString(c.Summary))
	if totalTickets == 0 {
		fmt.Println("There are no tickets in the campaign. Run \"campaigner add --help\" to find out how to add them.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "      %d\t-\ttotal tickets\t\n", totalTickets)
	fmt.Fprintf(w, "      %d\t%d%%\tpublished in Jira\t\n", totalPublishedJira, totalPublishedJira*100/totalTickets)
	fmt.Fprintf(w, "      %d\t%d%%\tpublished in Github\t\n", totalPublishedGithub, totalPublishedGithub*100/totalTickets)
	fmt.Fprintf(w, "      %d\t%d%%\tassigned\t\n", totalAssigned, totalAssigned*100/totalTickets)
	fmt.Fprintf(w, "      %d\t%d%%\tclosed\t\n\n", totalClosed, totalClosed*100/totalTickets)
	w.Flush()
}

func (c *Campaign) PrintList(publishedOnly, printLinks bool) {
	for _, t := range c.Tickets {
		if t.IsPublishedJira() {
			jiraLink := t.JiraLink
			if printLinks {
				jiraLink = c.GetJiraUrl(t)
			}

			var str string
			if t.IsPublishedGithub() {
				githubLink := fmt.Sprintf("#%d", t.GithubLink)
				if printLinks {
					githubLink = c.GetGithubUrl(t)
				}

				str = fmt.Sprintf("[%s / %s] %s", color.BlueString(jiraLink), color.CyanString(githubLink), t.Summary)
			} else {
				str = fmt.Sprintf("[%s] %s", color.BlueString(jiraLink), t.Summary)
			}
			if t.GithubStatus != "" {
				if t.IsClosed() {
					str += fmt.Sprintf(" (%s)", color.MagentaString(t.GithubStatus))
				} else {
					str += fmt.Sprintf(" (%s)", color.GreenString(t.GithubStatus))
				}
			}
			fmt.Println(str)
		} else if !publishedOnly {
			b, _ := json.Marshal(t)
			fmt.Printf("unpublished: %s\n", color.YellowString(string(b)))
		}
	}
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

func (c *Campaign) GetPublishedGithubTickets() []*Ticket {
	publishedTickets := []*Ticket{}
	for _, ticket := range c.Tickets {
		if ticket.IsPublishedGithub() {
			publishedTickets = append(publishedTickets, ticket)
		}
	}
	return publishedTickets
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
