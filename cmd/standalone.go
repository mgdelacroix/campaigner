package cmd

import (
	"fmt"
	"strings"

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

	cmd.Flags().StringP("username", "u", "", "the jira username")
	_ = cmd.MarkFlagRequired("username")
	cmd.Flags().StringP("token", "t", "", "the jira token")
	_ = cmd.MarkFlagRequired("token")
	cmd.Flags().StringP("summary", "s", "", "the summary of the ticket")
	_ = cmd.MarkFlagRequired("summary")
	cmd.Flags().StringP("template", "m", "", "the template to render the description of the ticket")
	_ = cmd.MarkFlagRequired("template")
	cmd.Flags().StringSliceP("vars", "v", []string{}, "the variables to use in the template")

	return cmd
}

func getVarMap(vars []string) (map[string]string, error) {
	varMap := map[string]string{}
	for _, v := range vars {
		parts := strings.Split(v, "=")
		if len(parts) < 2 {
			return nil, fmt.Errorf("cannot parse var %s", v)
		}
		varMap[parts[0]] = strings.Join(parts[1:], "")
	}
	return varMap, nil
}

func createJiraTicketStandaloneCmdF(cmd *cobra.Command, _ []string) error {
	username, _ := cmd.Flags().GetString("username")
	token, _ := cmd.Flags().GetString("token")
	summary, _ := cmd.Flags().GetString("summary")
	template, _ := cmd.Flags().GetString("template")
	vars, _ := cmd.Flags().GetStringSlice("vars")
	
	varMap, err := getVarMap(vars)
	if err != nil {
		return fmt.Errorf("error processing vars: %w")
	}
	
	// process template
	description := TBD()
	
	jiraClient, err := jira.NewClient(username, token)
	if err != nil {
		ErrorAndExit(cmd, err)
	}
	
	ticketKey, err := jiraClient.CreateTicket(summary, description)
	if err != nil {
		ErrorAndExit(cmd, err)
	}
	
	cmd.Printf("Ticket %s successfully created in JIRA")
	return nil
}
