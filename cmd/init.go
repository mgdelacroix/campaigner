package cmd

import (
	"git.ctrlz.es/mgdelacroix/campaigner/campaign"
	"git.ctrlz.es/mgdelacroix/campaigner/model"

	"github.com/spf13/cobra"
)

func InitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates a new campaign in the current directory",
		Args:  cobra.NoArgs,
		Run:   initCmdF,
	}

	cmd.Flags().StringP("url", "u", "", "the jira server URL")
	_ = cmd.MarkFlagRequired("url")
	cmd.Flags().StringP("project", "p", "", "the jira project key to associate the tickets with")
	_ = cmd.MarkFlagRequired("project")
	cmd.Flags().StringP("epic", "e", "", "the epic id to associate this campaign with")
	_ = cmd.MarkFlagRequired("epic")
	cmd.Flags().StringP("summary", "s", "", "the summary of the tickets")
	_ = cmd.MarkFlagRequired("summary")

	return cmd
}

func initCmdF(cmd *cobra.Command, _ []string) {
	url, _ := cmd.Flags().GetString("url")
	project, _ := cmd.Flags().GetString("project")
	epic, _ := cmd.Flags().GetString("epic")
	summary, _ := cmd.Flags().GetString("summary")

	cmp := &model.Campaign{
		Url:     url,
		Project: project,
		Epic:    epic,
		Summary: summary,
	}
	if err := campaign.Save(cmp); err != nil {
		ErrorAndExit(cmd, err)
	}
}
