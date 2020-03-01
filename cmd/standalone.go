package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"git.ctrlz.es/mgdelacroix/campaigner/jira"

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

	cmd.Flags().String("epic", "", "the jira epic id to associate the ticket with")
	_ = cmd.MarkFlagRequired("epic")
	cmd.Flags().String("team", "", "the team for the new ticket")
	_ = cmd.MarkFlagRequired("epic")
	cmd.Flags().String("username", "", "the jira username")
	_ = cmd.MarkFlagRequired("username")
	cmd.Flags().String("token", "", "the jira token")
	_ = cmd.MarkFlagRequired("token")
	cmd.Flags().String("summary", "", "the summary of the ticket")
	_ = cmd.MarkFlagRequired("summary")
	cmd.Flags().String("template", "", "the template to render the description of the ticket")
	_ = cmd.MarkFlagRequired("template")
	cmd.Flags().StringSliceP("vars", "v", []string{}, "the variables to use in the template")

	return cmd
}

func GetJiraTicketStandaloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-jira-ticket",
		Short: "Gets the ticket from jira",
		Args:  cobra.ExactArgs(1),
		Run:   getJiraTicketStandaloneCmdF,
	}

	cmd.Flags().String("username", "", "the jira username")
	_ = cmd.MarkFlagRequired("username")
	cmd.Flags().String("token", "", "the jira token")
	_ = cmd.MarkFlagRequired("token")

	return cmd
}

func getVarMap(vars []string) (map[string]string, error) {
	varMap := map[string]string{}
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
	epicId, _ := cmd.Flags().GetString("epic")
	team, _ := cmd.Flags().GetString("team")
	username, _ := cmd.Flags().GetString("username")
	token, _ := cmd.Flags().GetString("token")
	summaryTmplStr, _ := cmd.Flags().GetString("summary")
	templatePath, _ := cmd.Flags().GetString("template")
	vars, _ := cmd.Flags().GetStringSlice("vars")

	varMap, err := getVarMap(vars)
	if err != nil {
		return fmt.Errorf("error processing vars: %w", err)
	}

	sumTmpl, err := template.New("").Parse(summaryTmplStr)
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	var summaryBytes bytes.Buffer
	if err := sumTmpl.Execute(&summaryBytes, varMap); err != nil {
		ErrorAndExit(cmd, err)
	}
	summary := summaryBytes.String()

	descTmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	var descriptionBytes bytes.Buffer
	if err := descTmpl.Execute(&descriptionBytes, varMap); err != nil {
		ErrorAndExit(cmd, err)
	}
	description := descriptionBytes.String()

	jiraClient := jira.NewClient(username, token)

	ticketKey, err := jiraClient.CreateIssue(epicId, team, summary, description)
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	cmd.Printf("Ticket %s successfully created in JIRA", ticketKey)
	return nil
}

func getJiraTicketStandaloneCmdF(cmd *cobra.Command, args []string) {
	username, _ := cmd.Flags().GetString("username")
	token, _ := cmd.Flags().GetString("token")

	jiraClient := jira.NewClient(username, token)

	issue, err := jiraClient.GetIssue(args[0])
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	fmt.Printf("Key: %s\nStatus: %s\n", issue.Key, issue.Fields.Status.Name)
}
