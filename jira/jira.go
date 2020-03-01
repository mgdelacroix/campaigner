package jira

import (
	"encoding/json"
	"io"

	"git.ctrlz.es/mgdelacroix/campaigner/http"
)

type JiraClient struct {
	Username string
	Token    string
	Url      string
}

type JiraIssueFieldsStatus struct {
	Name string `json:"name"`
}

type JiraIssueFields struct {
	Status  JiraIssueFieldsStatus `json:"status"`
	Summary string                `json:"summary"`
}

type JiraIssue struct {
	Key    string          `json:"key"`
	Fields JiraIssueFields `json:"fields"`
}

func IssueFromJson(body io.Reader) (*JiraIssue, error) {
	var issue JiraIssue
	if err := json.NewDecoder(body).Decode(&issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

func NewClient(username, token string) *JiraClient {
	return &JiraClient{
		Username: username,
		Token:    token,
		Url:      "https://mattermost.atlassian.net/rest/api/2/",
	}
}

func (c *JiraClient) CreateIssue(epicId, team, summary, description string) (string, error) {
	data := map[string]interface{}{
		"fields": map[string]interface{}{
			"project":           map[string]interface{}{"key": "MM"},
			"summary":           summary,
			"description":       description,
			"issuetype":         map[string]interface{}{"name": "Story"},
			"customfield_10007": epicId,
			"customfield_11101": map[string]interface{}{"value": team},
		},
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	res, err := http.DoPost(c.Username, c.Token, c.Url+"issue/", body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	issue, err := IssueFromJson(res.Body)
	if err != nil {
		return "", err
	}

	return issue.Key, nil
}

func (c *JiraClient) GetIssue(issueNo string) (*JiraIssue, error) {
	res, err := http.DoGet(c.Username, c.Token, c.Url+"issue/"+issueNo)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	issue, err := IssueFromJson(res.Body)
	if err != nil {
		return nil, err
	}

	return issue, nil
}
