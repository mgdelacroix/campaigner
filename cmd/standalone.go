package cmd

import (
	"fmt"
	"strings"

	"git.ctrlz.es/mgdelacroix/campaigner/config"
	"git.ctrlz.es/mgdelacroix/campaigner/jira"
	"git.ctrlz.es/mgdelacroix/campaigner/model"

	"github.com/spf13/cobra"
)

func StandaloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "standalone",
		Short: "Standalone fire-and-forget commands",
	}

	cmd.AddCommand(
		CreateJiraTicketStandaloneCmd(),
		GetJiraTicketStandaloneCmd(),
	)

	return cmd
}

func CreateJiraTicketStandaloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-jira-ticket",
		Short: "Creates a jira ticket from a template",
		Args:  cobra.NoArgs,
		RunE:  createJiraTicketStandaloneCmdF,
	}

	cmd.Flags().String("url", "", "The jira server URL")
	_ = cmd.MarkFlagRequired("url")
	cmd.Flags().String("epic", "", "The jira epic id to associate the ticket with")
	_ = cmd.MarkFlagRequired("epic")
	cmd.Flags().StringP("project", "p", "", "The jira project key to associate the tickets with")
	_ = cmd.MarkFlagRequired("project")
	cmd.Flags().String("summary", "", "The summary of the ticket")
	_ = cmd.MarkFlagRequired("summary")
	cmd.Flags().String("template", "", "The template to render the description of the ticket")
	_ = cmd.MarkFlagRequired("template")
	cmd.Flags().String("username", "", "The jira username")
	cmd.Flags().String("token", "", "The jira token")
	cmd.Flags().StringSliceP("vars", "v", []string{}, "The variables to use in the template")
	cmd.Flags().Bool("dry-run", false, "Print the ticket information instead of creating it")

	return cmd
}

func GetJiraTicketStandaloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-jira-ticket",
		Short: "Gets the ticket from jira",
		Args:  cobra.ExactArgs(1),
		Run:   getJiraTicketStandaloneCmdF,
	}

	cmd.Flags().String("url", "", "The jira server URL")
	_ = cmd.MarkFlagRequired("url")
	cmd.Flags().String("username", "", "The jira username")
	cmd.Flags().String("token", "", "The jira token")

	return cmd
}

func getVarMap(vars []string) (map[string]interface{}, error) {
	varMap := map[string]interface{}{}
	for _, v := range vars {
		parts := strings.Split(v, "=")
		if len(parts) < 2 {
			return nil, fmt.Errorf("cannot parse var %s", v)
		}
		varMap[parts[0]] = strings.Join(parts[1:], "=")
	}
	return varMap, nil
}

func createJiraTicketStandaloneCmdF(cmd *cobra.Command, _ []string) error {
	url, _ := cmd.Flags().GetString("url")
	epic, _ := cmd.Flags().GetString("epic")
	project, _ := cmd.Flags().GetString("project")
	username, _ := cmd.Flags().GetString("username")
	token, _ := cmd.Flags().GetString("token")
	summary, _ := cmd.Flags().GetString("summary")
	template, _ := cmd.Flags().GetString("template")
	vars, _ := cmd.Flags().GetStringSlice("vars")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if username == "" || token == "" {
		cfg, err := config.ReadConfig()
		if err != nil {
			ErrorAndExit(cmd, err)
		}

		if username == "" {
			username = cfg.JiraUsername
		}
		if token == "" {
			token = cfg.JiraToken
		}
	}

	varMap, err := getVarMap(vars)
	if err != nil {
		return fmt.Errorf("error processing vars: %w", err)
	}

	jiraClient, err := jira.NewClient(url, username, token)
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	campaign := &model.Campaign{
		Epic:     epic,
		Project:  project,
		Summary:  summary,
		Template: template,
	}
	ticket := &model.Ticket{Data: varMap}

	issue, err := jiraClient.PublishTicket(ticket, campaign, dryRun)
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	cmd.Printf("Ticket %s successfully created in JIRA", issue.Key)
	return nil
}

func getJiraTicketStandaloneCmdF(cmd *cobra.Command, args []string) {
	url, _ := cmd.Flags().GetString("url")
	username, _ := cmd.Flags().GetString("username")
	token, _ := cmd.Flags().GetString("token")

	if username == "" || token == "" {
		cfg, err := config.ReadConfig()
		if err != nil {
			ErrorAndExit(cmd, err)
		}

		if username == "" {
			username = cfg.JiraUsername
		}
		if token == "" {
			token = cfg.JiraToken
		}
	}

	jiraClient, err := jira.NewClient(url, username, token)
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	issue, err := jiraClient.GetIssue(args[0])
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	fmt.Printf("Summary: %s\nKey: %s\nStatus: %s\nAsignee: %s\n", issue.Fields.Summary, issue.Key, issue.Fields.Status.Name, issue.Fields.Assignee.DisplayName)
}
