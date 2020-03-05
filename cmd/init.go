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

	cmd.Flags().StringP("url", "u", "", "The jira server URL")
	_ = cmd.MarkFlagRequired("url")
	cmd.Flags().StringP("project", "p", "", "The jira project key to associate the tickets with")
	_ = cmd.MarkFlagRequired("project")
	cmd.Flags().StringP("epic", "e", "", "The epic id to associate this campaign with")
	_ = cmd.MarkFlagRequired("epic")
	cmd.Flags().StringP("summary", "s", "", "The summary of the tickets")
	_ = cmd.MarkFlagRequired("summary")
	cmd.Flags().StringP("template", "t", "", "The template path for the description of the tickets")
	_ = cmd.MarkFlagRequired("template")

	return cmd
}

func initCmdF(cmd *cobra.Command, _ []string) {
	url, _ := cmd.Flags().GetString("url")
	project, _ := cmd.Flags().GetString("project")
	epic, _ := cmd.Flags().GetString("epic")
	summary, _ := cmd.Flags().GetString("summary")
	template, _ := cmd.Flags().GetString("template")

	cmp := &model.Campaign{
		Url:      url,
		Project:  project,
		Epic:     epic,
		Summary:  summary,
		Template: template,
	}
	if err := campaign.Save(cmp); err != nil {
		ErrorAndExit(cmd, err)
	}
}
