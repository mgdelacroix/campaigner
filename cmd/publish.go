package cmd

import (
	"fmt"

	"github.com/mgdelacroix/campaigner/app"
	"github.com/mgdelacroix/campaigner/model"

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

func TicketPublishCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ticket",
		Short: "Publishes an already existing jira ticket in Github",
		Long: `This command publishes a ticket that is outside the campaign and already created in jira, into Github.

It is intended to use for standalone Help Wanted tickets that don't fit into a campaing, but nonetheless need to be published in Github and linked back. It does require of a campaign.json file that describes the connection with both the Jira instance and the GitHub repository, but it will never modify it, so the command can be run pointing to a previously existing campaign which connection details match with the ones that apply for the ticket.

Github labels will not be read from the campaign.json file, so they need to be specified with the --label flag if wanted.`,
		Example: `  # if we don't want any github label to be added to the ticket
  $ campaigner publish ticket MM-1234

  # if we want to add some labels in github
  $ campaigner publish ticket MM-1234 --label Tech/Go --label "Help Wanted"

  # if we want to use a campaign file outside the current directory
  $ campaigner publish ticket MM-1234 --campaign ~/campaigns/standalone.json`,
		Args: cobra.ExactArgs(1),
		RunE: withAppE(ticketPublishCmdF),
	}

	cmd.Flags().StringSliceP("label", "l", []string{}, "the labels to add to the Github issues")
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
		TicketPublishCmd(),
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

func ticketPublishCmdF(a *app.App, cmd *cobra.Command, args []string) error {
	labels, _ := cmd.Flags().GetStringSlice("label")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	jiraTicketId := args[0]

	jiraIssue, err := a.GetIssue(jiraTicketId)
	if err != nil {
		ErrorAndExit(cmd, fmt.Errorf("cannot get jira issue %q: %w", jiraTicketId, err))
	}

	ticket := &model.Ticket{
		Summary:     jiraIssue.Fields.Summary,
		Description: jiraIssue.Fields.Description,
		JiraLink:    jiraTicketId,
	}
	// update the campaign labels only to publish the ticket
	a.Campaign.Github.Labels = labels

	githubIssue, err := a.PublishInGithub(ticket, dryRun)
	if err != nil {
		ErrorAndExit(cmd, fmt.Errorf("cannot publish ticket %q in github: %w", jiraTicketId, err))
	}

	if dryRun {
		return nil
	}

	ticket.GithubLink = githubIssue.GetNumber()
	ticket.GithubStatus = githubIssue.GetState()

	cmd.Printf("Issue published: https://github.com/%s/issues/%d\n", a.Campaign.Github.Repo, ticket.GithubLink)

	if err := a.UpdateJiraAfterGithub(ticket); err != nil {
		ErrorAndExit(cmd, fmt.Errorf("error updating Jira info for %q after publishing in Github", jiraTicketId))
	}

	return nil
}
