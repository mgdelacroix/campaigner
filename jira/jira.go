package jira

import (
	"bytes"
	"fmt"
	"text/template"

	"git.ctrlz.es/mgdelacroix/campaigner/campaign"
	"git.ctrlz.es/mgdelacroix/campaigner/model"

	jira "gopkg.in/andygrunwald/go-jira.v1"
)

type JiraClient struct {
	*jira.Client
}

func NewClient(url, username, token string) (*JiraClient, error) {
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: token,
	}

	client, err := jira.NewClient(tp.Client(), url)
	if err != nil {
		return nil, err
	}

	return &JiraClient{client}, nil
}

func (c *JiraClient) GetIssueFromTicket(ticket *model.Ticket, cmp *model.Campaign) (*jira.Issue, error) {
	summaryTmpl, err := template.New("").Parse(cmp.Summary)
	if err != nil {
		return nil, err
	}

	var summaryBytes bytes.Buffer
	if err := summaryTmpl.Execute(&summaryBytes, ticket.Data); err != nil {
		return nil, err
	}
	summary := summaryBytes.String()

	descriptionTemplate, err := template.ParseFiles(cmp.Template)
	if err != nil {
		return nil, err
	}

	var descriptionBytes bytes.Buffer
	if err := descriptionTemplate.Execute(&descriptionBytes, ticket.Data); err != nil {
		return nil, err
	}
	description := descriptionBytes.String()

	data := map[string]string{
		"Description": description,
		"Summary":     summary,
		"Project":     cmp.Project,
		"Issue Type":  cmp.IssueType,
		"Epic Link":   cmp.Epic,
	}

	if team, ok := ticket.Data["team"]; ok {
		data["team"] = team.(string)
	}

	createMetaInfo, _, err := c.Issue.GetCreateMeta(cmp.Project)
	if err != nil {
		return nil, err
	}

	project := createMetaInfo.GetProjectWithKey(cmp.Project)
	if project == nil {
		return nil, fmt.Errorf("Error retrieving project with key %s", cmp.Project)
	}

	issueType := project.GetIssueTypeWithName(cmp.IssueType)
	if issueType == nil {
		return nil, fmt.Errorf("Error retrieving issue type with name Story")
	}

	issue, err := jira.InitIssueWithMetaAndFields(project, issueType, data)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func (c *JiraClient) PublishTicket(ticket *model.Ticket, cmp *model.Campaign) (*jira.Issue, error) {
	issue, err := c.GetIssueFromTicket(ticket, cmp)
	if err != nil {
		return nil, err
	}

	newIssue, _, err := c.Issue.Create(issue)
	if err != nil {
		return nil, err
	}

	return newIssue, nil
}

func (c *JiraClient) GetIssue(issueNo string) (*jira.Issue, error) {
	issue, _, err := c.Issue.Get(issueNo, nil)
	if err != nil {
		return nil, err
	}
	return issue, nil
}

func (c *JiraClient) PublishNextTicket(cmp *model.Campaign) (bool, error) {
	ticket := cmp.NextUnpublishedTicket()
	if ticket == nil {
		return false, nil
	}

	issue, err := c.PublishTicket(ticket, cmp)
	if err != nil {
		return false, err
	}

	ticket.JiraLink = issue.Key
	if err := campaign.Save(cmp); err != nil {
		return false, err
	}
	return true, nil
}

func (c *JiraClient) PublishAll(cmp *model.Campaign) (int, error) {
	count := 0
	for {
		next, err := c.PublishNextTicket(cmp)
		if err != nil {
			return count, err
		}
		if !next {
			break
		}
		count++
	}
	return count, nil
}

func (c *JiraClient) PublishBatch(cmp *model.Campaign, batch int) error {
	return nil
}
