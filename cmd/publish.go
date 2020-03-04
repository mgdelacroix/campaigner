package cmd

import (
	"github.com/spf13/cobra"
)

func JiraPublishCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "jira",
		Short: "Publishes the campaign tickets in JIRA",
		Args:  cobra.NoArgs,
		Run:   jiraPublishCmdF,
	}
}

func PublishCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Publishes the campaign tickets in different providers",
	}

	cmd.AddCommand(
		JiraPublishCmd(),
	)

	return cmd
}

func jiraPublishCmdF(_ *cobra.Command, _ []string) {}
