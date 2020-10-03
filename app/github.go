package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"git.ctrlz.es/mgdelacroix/campaigner/model"

	"github.com/StevenACoffman/j2m"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func (a *App) InitGithubClient() error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: a.Campaign.Github.Token})
	tc := oauth2.NewClient(ctx, ts)

	a.GithubClient = github.NewClient(tc)
	return nil
}

func getFooterTemplate(ticket *model.Ticket, templatePath string) (string, error) {
	footerTmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var footerBytes bytes.Buffer
	if err := footerTmpl.Execute(&footerBytes, ticket); err != nil {
		return "", err
	}
	return footerBytes.String(), nil
}

func (a *App) PublishInGithub(ticket *model.Ticket, dryRun bool) (*github.Issue, error) {
	mdDescription := j2m.JiraToMD(ticket.Description)
	if a.Campaign.FooterTemplate != "" {
		footer, err := getFooterTemplate(ticket, a.Campaign.FooterTemplate)
		if err != nil {
			return nil, err
		}

		mdDescription += "\n" + footer
	}

	issueRequest := &github.IssueRequest{
		Title:  &ticket.Summary,
		Body:   &mdDescription,
		Labels: &a.Campaign.Github.Labels,
	}

	if dryRun {
		b, _ := json.MarshalIndent(issueRequest, "", "  ")
		fmt.Println(string(b))
		return &github.Issue{
			Title: issueRequest.Title,
			Body:  issueRequest.Body,
		}, nil
	}

	owner, repo := a.Campaign.RepoComponents()
	newIssue, _, err := a.GithubClient.Issues.Create(context.Background(), owner, repo, issueRequest)
	if err != nil {
		return nil, err
	}
	return newIssue, nil
}

func (a *App) PublishNextInGithub(dryRun bool) (bool, error) {
	ticket := a.Campaign.NextGithubUnpublishedTicket()
	if ticket == nil {
		return false, nil
	}

	issue, err := a.PublishInGithub(ticket, dryRun)
	if err != nil {
		return false, err
	}

	if dryRun {
		return true, nil
	}

	ticket.GithubLink = issue.GetNumber()
	ticket.GithubStatus = issue.GetState()
	if err := a.Save(); err != nil {
		return false, err
	}

	// ToDo: print here the newly created issue

	if !dryRun {
		if err := a.UpdateJiraAfterGithub(ticket); err != nil {
			fmt.Fprintf(os.Stderr, "error updating Jira info for %s after publishing in Github\n", ticket.JiraLink)
		}
	}

	return true, nil
}

func (a *App) PublishAllInGithub(dryRun bool) (int, error) {
	count := 0
	for {
		next, err := a.PublishNextInGithub(dryRun)
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

func (a *App) PublishBatchInGithub(batch int, dryRun bool) error {
	for i := 1; i <= batch; i++ {
		next, err := a.PublishNextInGithub(dryRun)
		if err != nil {
			return err
		}
		if !next {
			return nil
		}
	}
	return nil
}

func (a *App) GithubSync() error {
	tickets := a.Campaign.GetPublishedGithubTickets()
	total := len(tickets)
	owner, repo := a.Campaign.RepoComponents()

	for i, ticket := range tickets {
		fmt.Printf("\rUpdating ticket %d of %d", i+1, total)

		issue, _, err := a.GithubClient.Issues.Get(context.Background(), owner, repo, ticket.GithubLink)
		if err != nil {
			return err
		}

		assignee := issue.GetAssignee()
		if assignee != nil {
			ticket.GithubAssignee = assignee.GetLogin()
		}
		ticket.GithubStatus = issue.GetState()
	}
	fmt.Print("\n")

	return a.Save()
}

func (a *App) ListLabels() ([]string, error) {
	owner, repo := a.Campaign.RepoComponents()
	opts := &github.ListOptions{Page: 0, PerPage: 100}
	labels, _, err := a.GithubClient.Issues.ListLabels(context.Background(), owner, repo, opts)
	if err != nil {
		return nil, err
	}

	strLabels := make([]string, len(labels))
	for i, label := range labels {
		strLabels[i] = *label.Name
	}

	return strLabels, nil
}

func (a *App) CheckLabels(labels []string) (bool, []string, error) {
	ghLabels, err := a.ListLabels()
	if err != nil {
		return false, nil, err
	}

	badLabels := []string{}
	for _, label := range labels {
		exists := false
		for _, ghLabel := range ghLabels {
			if label == ghLabel {
				exists = true
			}
		}

		if !exists {
			badLabels = append(badLabels, label)
		}
	}

	if len(badLabels) == 0 {
		return true, nil, nil
	}
	return false, badLabels, nil
}
