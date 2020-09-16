package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
	cmd.Flags().StringP("github-issue", "g", "", "syncs a ticket by GitHub issue number")

	return cmd
}

func syncCmdF(cmd *cobra.Command, _ []string) {
	ErrorAndExit(cmd, fmt.Errorf("Not implemented yet"))
}
