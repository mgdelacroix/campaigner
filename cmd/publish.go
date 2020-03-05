package cmd

import (
	"github.com/spf13/cobra"
)

func JiraPublishCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jira",
		Short: "Publishes the campaign tickets in JIRA",
		Args:  cobra.NoArgs,
		RunE:  jiraPublishCmdF,
	}

	cmd.Flags().BoolP("all", "a", false, "Publish all the tickets of the campaign")
	cmd.Flags().IntP("batch", "b", 0, "Number of tickets to publish")

	return cmd
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

func jiraPublishCmdF(cmd *cobra.Command, _ []string) error {
	/*
		all, _ := cmd.Flags().GetBool("all")
		batch, _ := cmd.Flags().GetInt("batch")

		if !all && batch == 0 {
			return fmt.Errorf("One of --all or --batch flags is required")
		}
	*/

	return nil
}
