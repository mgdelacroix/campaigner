package app

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/go-github/v32/github"
	jira "gopkg.in/andygrunwald/go-jira.v1"

	"github.com/mgdelacroix/campaigner/model"
)

type App struct {
	Path string

	JiraClient   *jira.Client
	GithubClient *github.Client
	Campaign     *model.Campaign
}

func SaveCampaign(campaign *model.Campaign, path string) error {
	marshaledCampaign, err := json.MarshalIndent(campaign, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, marshaledCampaign, 0600); err != nil {
		return fmt.Errorf("cannot save campaign: %w", err)
	}
	return nil
}

func (a *App) Save() error {
	return SaveCampaign(a.Campaign, a.Path)
}

func (a *App) Load() error {
	if _, err := os.Stat("."); err != nil {
		return fmt.Errorf("cannot read campaign: %w", err)
	}

	fileContents, err := os.ReadFile(a.Path)
	if err != nil {
		return fmt.Errorf("there was a problem reading the campaign file: %w", err)
	}

	var campaign model.Campaign
	if err := json.Unmarshal(fileContents, &campaign); err != nil {
		return fmt.Errorf("there was a problem parsing the campaign file: %w", err)
	}

	a.Campaign = &campaign
	return nil
}

func (a *App) InitClients() error {
	if err := a.InitGithubClient(); err != nil {
		return err
	}
	if err := a.InitJiraClient(); err != nil {
		return err
	}
	return nil
}

func NewApp(path string) (*App, error) {
	app := &App{Path: path}

	if err := app.Load(); err != nil {
		return nil, err
	}
	if err := app.InitClients(); err != nil {
		return nil, err
	}
	return app, nil
}
