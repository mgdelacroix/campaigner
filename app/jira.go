package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"git.ctrlz.es/mgdelacroix/campaigner/model"

	jira "gopkg.in/andygrunwald/go-jira.v1"
)

const GithubIssueJiraField = "customfield_11106"

func (a *App) InitJiraClient() error {
	tp := jira.BasicAuthTransport{
		Username: a.Campaign.Jira.Username,
		Password: a.Campaign.Jira.Token,
	}

	client, err := jira.NewClient(tp.Client(), a.Campaign.Jira.Url)
	if err != nil {
		return err
	}

	a.JiraClient = client
	return nil
}

func (a *App) UpdateJiraAfterGithub(ticket *model.Ticket) error {
	data := map[string]interface{}{
		"fields": map[string]interface{}{
			GithubIssueJiraField: a.Campaign.GetGithubUrl(ticket),
			"fixVersions": []map[string]interface{}{
				{
					"name": "Help Wanted",
				},
			},
		},
	}

	_, err := a.JiraClient.Issue.UpdateIssue(ticket.JiraLink, data)
	return err
}

func (a *App) GetJiraIssueFromTicket(ticket *model.Ticket) (*jira.Issue, error) {
	summaryTmpl, err := template.New("").Parse(a.Campaign.Summary)
	if err != nil {
		return nil, err
	}

	var summaryBytes bytes.Buffer
	if err := summaryTmpl.Execute(&summaryBytes, ticket.Data); err != nil {
		return nil, err
	}
	summary := summaryBytes.String()

	descriptionTemplate, err := template.ParseFiles(a.Campaign.IssueTemplate)
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
		"Project":     a.Campaign.Jira.Project,
		"Issue Type":  a.Campaign.Jira.IssueType,
		"Epic Link":   a.Campaign.Jira.Epic,
	}

	createMetaInfo, _, err := a.JiraClient.Issue.GetCreateMeta(a.Campaign.Jira.Project)
	if err != nil {
		return nil, err
	}

	project := createMetaInfo.GetProjectWithKey(a.Campaign.Jira.Project)
	if project == nil {
		return nil, fmt.Errorf("Error retrieving project with key %s", a.Campaign.Jira.Project)
	}

	issueType := project.GetIssueTypeWithName(a.Campaign.Jira.IssueType)
	if issueType == nil {
		return nil, fmt.Errorf("Error retrieving issue type with name Story")
	}

	issue, err := jira.InitIssueWithMetaAndFields(project, issueType, data)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func (a *App) PublishInJira(ticket *model.Ticket, dryRun bool) (*jira.Issue, error) {
	issue, err := a.GetJiraIssueFromTicket(ticket)
	if err != nil {
		return nil, err
	}

	if dryRun {
		b, _ := json.MarshalIndent(issue, "", "  ")
		fmt.Println(string(b))
		return issue, nil
	}

	newIssue, _, err := a.JiraClient.Issue.Create(issue)
	if err != nil {
		return nil, err
	}

	return newIssue, nil
}

func (a *App) GetIssue(issueNo string) (*jira.Issue, error) {
	issue, _, err := a.JiraClient.Issue.Get(issueNo, nil)
	if err != nil {
		return nil, err
	}
	return issue, nil
}

func (a *App) PublishNextInJira(dryRun bool) (bool, error) {
	ticket := a.Campaign.NextJiraUnpublishedTicket()
	if ticket == nil {
		return false, nil
	}

	issue, err := a.PublishInJira(ticket, dryRun)
	if err != nil {
		return false, err
	}

	if dryRun {
		return true, nil
	}

	issue, _, err = a.JiraClient.Issue.Get(issue.Key, nil)
	if err != nil {
		return false, err
	}

	ticket.JiraLink = issue.Key
	ticket.Summary = issue.Fields.Summary
	ticket.Description = issue.Fields.Description
	ticket.JiraStatus = issue.Fields.Status.Name
	if err := a.Save(); err != nil {
		return false, err
	}
	return true, nil
}

func (a *App) PublishAllInJira(dryRun bool) (int, error) {
	count := 0
	for {
		next, err := a.PublishNextInJira(dryRun)
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

func (a *App) PublishBatchInJira(batch int, dryRun bool) error {
	for i := 1; i <= batch; i++ {
		next, err := a.PublishNextInJira(dryRun)
		if err != nil {
			return err
		}
		if !next {
			return nil
		}
	}
	return nil
}

func (a *App) GetTicketsFromJiraEpic() ([]*model.Ticket, error) {
	jql := fmt.Sprintf("project = %s AND type = %s AND \"Epic Link\" = %s", a.Campaign.Jira.Project, a.Campaign.Jira.IssueType, a.Campaign.Jira.Epic)

	page := 0
	maxPerPage := 50
	issues := []jira.Issue{}
	for {
		opts := &jira.SearchOptions{StartAt: maxPerPage * page, MaxResults: maxPerPage}
		pageIssues, _, err := a.JiraClient.Issue.Search(jql, opts)
		if err != nil {
			return nil, err
		}

		issues = append(issues, pageIssues...)
		if len(pageIssues) < maxPerPage {
			break
		}
		page++
	}

	tickets := []*model.Ticket{}
	for _, issue := range issues {
		// ToDo: if they have github link, fill and fetch github data
		ticket := &model.Ticket{
			JiraLink:    issue.Key,
			JiraStatus:  issue.Fields.Status.Name,
			Summary:     issue.Fields.Summary,
			Description: issue.Fields.Description,
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}
