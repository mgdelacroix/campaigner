package jira

import (
	"encoding/json"
	"io/ioutil"

	"git.ctrlz.es/mgdelacroix/campaigner/http"
)

type JiraClient struct {
	Username string
	Token    string
	Url      string
}

type JiraIssue struct {
	Key string `json:"key"`
}

func NewClient(username, token string) *JiraClient {
	return &JiraClient{
		Username: username,
		Token:    token,
		Url:      "https://mattermost.atlassian.net/rest/api/2",
	}
}

func (c *JiraClient) CreateTicket(epicId, team, summary, description string) (string, error) {
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

	res, err := http.DoPost(c.Username, c.Token, c.Url+"/issue/", body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var issue JiraIssue
	if err := json.Unmarshal(respBody, &issue); err != nil {
		return "", err
	}

	return issue.Key, nil
}
