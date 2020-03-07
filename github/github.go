package github

import (
	"context"

	"git.ctrlz.es/mgdelacroix/campaigner/campaign"
	"git.ctrlz.es/mgdelacroix/campaigner/model"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	*github.Client
	Repo string
}

func NewClient(repo, token string) *GithubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &GithubClient{
		Client: client,
		Repo:   repo,
	}
}

func (c *GithubClient) PublishTicket(ticket *model.Ticket, cmp *model.Campaign, dryRun bool) (*github.Issue, error) {
	return nil, nil
}

func (c *GithubClient) PublishNextTicket(cmp *model.Campaign, dryRun bool) (bool, error) {
	ticket := cmp.NextGithubUnpublishedTicket()
	if ticket == nil {
		return false, nil
	}

	issue, err := c.PublishTicket(ticket, cmp, dryRun)
	if err != nil {
		return false, err
	}

	if dryRun {
		return true, nil
	}

	ticket.GithubLink = *issue.ID
	// move this to a publish service that can do both github and
	// jira, as we need to update a jira issue field with the github
	// link
	if err := campaign.Save(cmp); err != nil {
		return false, err
	}
	return true, nil
}

func (c *GithubClient) PublishAll(cmp *model.Campaign, dryRun bool) (int, error) {
	count := 0
	for {
		next, err := c.PublishNextTicket(cmp, dryRun)
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

func (c *GithubClient) PublishBatch(cmp *model.Campaign, batch int, dryRun bool) error {
	for i := 0; i <= batch; i++ {
		next, err := c.PublishNextTicket(cmp, dryRun)
		if err != nil {
			return err
		}
		if !next {
			return nil
		}
	}
	return nil
}
