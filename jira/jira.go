package jira

import (
	"bytes"
	"fmt"
	"text/template"

	"git.ctrlz.es/mgdelacroix/campaigner/model"

	jira "gopkg.in/andygrunwald/go-jira.v1"
)

type JiraClient struct {
	*jira.Client
}

func (c *JiraClient) GetIssueFromTicket(ticket *model.Ticket, campaign *model.Campaign) (*jira.Issue, error) {
	summaryTmpl, err := template.New("").Parse(campaign.Summary)
	if err != nil {
		return nil, err
	}

	var summaryBytes bytes.Buffer
	if err := summaryTmpl.Execute(&summaryBytes, ticket.Data); err != nil {
		return nil, err
	}
	summary := summaryBytes.String()

	descriptionTemplate, err := template.ParseFiles(campaign.Template)
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
		"Project":     campaign.Project,
		"Issue Type":  campaign.IssueType,
		"Epic Link":   campaign.Epic,
	}

	if team, ok := ticket.Data["team"]; ok {
		data["team"] = team.(string)
	}

	createMetaInfo, _, err := c.Issue.GetCreateMeta(campaign.Project)
	if err != nil {
		return nil, err
	}

	project := createMetaInfo.GetProjectWithKey(campaign.Project)
	if project == nil {
		return nil, fmt.Errorf("Error retrieving project with key %s", campaign.Project)
	}

	issueType := project.GetIssueTypeWithName(campaign.IssueType)
	if issueType == nil {
		return nil, fmt.Errorf("Error retrieving issue type with name Story")
	}

	issue, err := jira.InitIssueWithMetaAndFields(project, issueType, data)
	if err != nil {
		return nil, err
	}

	return issue, nil
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

func (c *JiraClient) CreateTicket(ticket *model.Ticket, campaign *model.Campaign) (*jira.Issue, error) {
	issue, err := c.GetIssueFromTicket(ticket, campaign)
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
