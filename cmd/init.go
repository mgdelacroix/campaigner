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

	cmd.Flags().StringP("epic", "e", "", "the epic id to associate this campaign with")
	_ = cmd.MarkFlagRequired("epic")
	cmd.Flags().StringP("summary", "s", "", "the summary of the tickets. Can contain the variables {{.Filename}}, {{.LineNo}} and {{.Text}}")
	_ = cmd.MarkFlagRequired("summary")

	return cmd
}

func initCmdF(cmd *cobra.Command, _ []string) {
	epic, _ := cmd.Flags().GetString("epic")
	summary, _ := cmd.Flags().GetString("summary")
	if err := campaign.Save(&model.Campaign{Epic: epic, Summary: summary}); err != nil {
		ErrorAndExit(cmd, err)
	}
}
