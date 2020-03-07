package github

import (
	"context"

	"golang.org/x/oauth2"
	"github.com/google/go-github/v29/github"
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
		Repo: repo,
	}
}
