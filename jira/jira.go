package jira

import (
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

func (c *JiraClient) CreateIssue(epicId, team, summary, description string) (string, error) {
	/*
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
	*/
	return "", nil
}

func (c *JiraClient) GetIssue(issueNo string) (*jira.Issue, error) {
	issue, _, err := c.Issue.Get(issueNo, nil)
	if err != nil {
		return nil, err
	}
	return issue, nil
}
