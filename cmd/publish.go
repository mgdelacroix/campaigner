package cmd

import (
	"fmt"

	"git.ctrlz.es/mgdelacroix/campaigner/app"

	"github.com/spf13/cobra"
)

func JiraPublishCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jira",
		Short: "Publishes the campaign tickets in jira",
		Args:  cobra.NoArgs,
		RunE:  withAppE(jiraPublishCmdF),
	}

	cmd.Flags().BoolP("all", "a", false, "Publish all the tickets of the campaign")
	cmd.Flags().IntP("batch", "b", 0, "Number of tickets to publish")
	cmd.Flags().Bool("dry-run", false, "Print the tickets information instead of publishing them")

	return cmd
}

func GithubPublishCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "github",
		Short: "Publishes the campaign tickets in github",
		Args:  cobra.NoArgs,
		RunE:  withAppE(githubPublishCmdF),
	}

	cmd.Flags().BoolP("all", "a", false, "Publish all the tickets of the campaign")
	cmd.Flags().IntP("batch", "b", 0, "Number of tickets to publish")
	cmd.Flags().Bool("dry-run", false, "Print the tickets information instead of publishing them")

	return cmd
}

func PublishCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Publishes campaign tickets",
		Long:  "Publishes the campaign tickets in both Jira and Github",
	}

	cmd.AddCommand(
		GithubPublishCmd(),
		JiraPublishCmd(),
	)

	return cmd
}

func jiraPublishCmdF(a *app.App, cmd *cobra.Command, _ []string) error {
	all, _ := cmd.Flags().GetBool("all")
	batch, _ := cmd.Flags().GetInt("batch")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if !all && batch == 0 {
		return fmt.Errorf("One of --all or --batch flags is required")
	}

	if all {
		count, err := a.PublishAllInJira(cmd.OutOrStdout(), dryRun)
		if err != nil {
			ErrorAndExit(cmd, err)
		}
		cmd.Printf("\nAll %d tickets successfully published in jira\n", count)
	} else {
		if err := a.PublishBatchInJira(cmd.OutOrStdout(), batch, dryRun); err != nil {
			ErrorAndExit(cmd, err)
		}
		cmd.Printf("\nBatch of %d tickets successfully published in jira\n", batch)
	}

	return nil
}

func githubPublishCmdF(a *app.App, cmd *cobra.Command, _ []string) error {
	all, _ := cmd.Flags().GetBool("all")
	batch, _ := cmd.Flags().GetInt("batch")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if !all && batch == 0 {
		return fmt.Errorf("One of --all or --batch flags is required")
	}

	if all {
		count, err := a.PublishAllInGithub(cmd.OutOrStdout(), dryRun)
		if err != nil {
			ErrorAndExit(cmd, err)
		}
		cmd.Printf("\nAll %d tickets successfully published in github\n", count)
	} else {
		if err := a.PublishBatchInGithub(cmd.OutOrStdout(), batch, dryRun); err != nil {
			ErrorAndExit(cmd, err)
		}
		cmd.Printf("\nBatch of %d tickets successfully published in github\n", batch)
	}

	return nil
}
