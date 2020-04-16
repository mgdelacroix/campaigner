package cmd

import (
	"strings"

	"git.ctrlz.es/mgdelacroix/campaigner/campaign"
	"git.ctrlz.es/mgdelacroix/campaigner/model"

	"github.com/spf13/cobra"
)

func InitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates a new campaign in the current directory",
		Example: `  campaigner init \
    --url http://my-jira-instance.com \
    --epic ASD-27 \
    --issue-type Story \
    --summary 'Refactor {{.function}} to inject the configuration service' \
    --template ./refactor-config.tmpl`,
		Args:  cobra.NoArgs,
		Run:   initCmdF,
	}

	cmd.Flags().StringP("url", "u", "", "The jira server URL")
	_ = cmd.MarkFlagRequired("url")
	cmd.Flags().StringP("epic", "e", "", "The epic id to associate this campaign with")
	_ = cmd.MarkFlagRequired("epic")
	cmd.Flags().StringP("summary", "s", "", "The summary of the tickets")
	_ = cmd.MarkFlagRequired("summary")
	cmd.Flags().StringP("template", "t", "", "The template path for the description of the tickets")
	_ = cmd.MarkFlagRequired("template")
	cmd.Flags().StringP("issue-type", "i", "Story", "The issue type to create the tickets as")

	return cmd
}

func initCmdF(cmd *cobra.Command, _ []string) {
	url, _ := cmd.Flags().GetString("url")
	epic, _ := cmd.Flags().GetString("epic")
	summary, _ := cmd.Flags().GetString("summary")
	template, _ := cmd.Flags().GetString("template")
	issueType, _ := cmd.Flags().GetString("issue-type")

	project := strings.Split(epic, "-")[0]

	cmp := &model.Campaign{
		Url:       url,
		Project:   project,
		Epic:      epic,
		IssueType: issueType,
		Summary:   summary,
		Template:  template,
	}
	if err := campaign.Save(cmp); err != nil {
		ErrorAndExit(cmd, err)
	}
}
