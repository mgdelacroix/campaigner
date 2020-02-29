package model

type Config struct {
	GithubToken  string `json:"github_token"`
	JiraUsername string `json:"jira_username"`
	JiraToken    string `json:"jira_token"`
}
