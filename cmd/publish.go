package cmd

import (
	"fmt"

	"git.ctrlz.es/mgdelacroix/campaigner/campaign"
	"git.ctrlz.es/mgdelacroix/campaigner/config"
	"git.ctrlz.es/mgdelacroix/campaigner/jira"

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
	cmd.Flags().Bool("dry-run", false, "Print the tickets information instead of publishing them")

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
	all, _ := cmd.Flags().GetBool("all")
	batch, _ := cmd.Flags().GetInt("batch")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if !all && batch == 0 {
		return fmt.Errorf("One of --all or --batch flags is required")
	}

	cfg, err := config.ReadConfig()
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	cmp, err := campaign.Read()
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	jiraClient, err := jira.NewClient(cmp.Url, cfg.JiraUsername, cfg.JiraToken)
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	if all {
		count, err := jiraClient.PublishAll(cmp, dryRun)
		if err != nil {
			ErrorAndExit(cmd, err)
		}
		cmd.Printf("All %d tickets successfully published in jira\n", count)
	} else {
		if err := jiraClient.PublishBatch(cmp, batch, dryRun); err != nil {
			ErrorAndExit(cmd, err)
		}
		cmd.Printf("Batch of %d tickets successfully published in jira\n", batch)
	}

	return nil
}