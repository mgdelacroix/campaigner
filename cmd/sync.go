package cmd

import (
	"fmt"
	"encoding/json"

	"github.com/spf13/cobra"

	"git.ctrlz.es/mgdelacroix/campaigner/campaign"
	"git.ctrlz.es/mgdelacroix/campaigner/jira"
	"git.ctrlz.es/mgdelacroix/campaigner/model"
)

func SyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Syncs the tickets",
		Long:  "Synchronizes the status of the published tickets with remote providers",
		Args:  cobra.NoArgs,
		Run:   syncCmdF,
	}

	cmd.Flags().BoolP("all", "a", false, "syncs all the published tickets")
	cmd.Flags().StringP("jira-issue", "j", "", "syncs a ticket by Jira issue number")
	cmd.Flags().IntP("github-issue", "g", 0, "syncs a ticket by GitHub issue number")

	return cmd
}

func syncCmdF(cmd *cobra.Command, _ []string) {
	jiraIssue, _ := cmd.Flags().GetString("jira-issue")
	// githubIssue, _ := cmd.Flags().GetInt()

	// check that one is defined, or all

	cmp, err := campaign.Read()
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	var ticket *model.Ticket
	if jiraIssue != "" {
		ticket = cmp.GetByJiraIssue(jiraIssue)
		if ticket == nil {
			ErrorAndExit(cmd, fmt.Errorf("Could not find jira issue %s", jiraIssue))
		}
	}

	jiraClient, err := jira.NewClient(cmp.Jira.Url, cmp.Jira.Username, cmp.Jira.Token)
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	i, _, err := jiraClient.Issue.Get(ticket.JiraIssue(), nil)
	if err != nil {
		ErrorAndExit(cmd, err)
	}
	b, _ := json.MarshalIndent(i.Fields, "", "  ")

	fmt.Printf(string(b))
}
