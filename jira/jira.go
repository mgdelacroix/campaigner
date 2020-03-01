package jira

import (
	"fmt"
)

type JiraClient struct {
	Username string
	Token string
}

func NewClient(username, token string) *JiraClient {
	return &JiraClient{
		Username: username,
		Token: token,
	}
}

func (c *JiraClient) CreateTicket(summary, description string) (string, error) {
	fmt.Printf("Summary: %s\nDescription: %s\n", summary, description)
	return "", nil
}
