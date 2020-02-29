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

	return cmd
}

func initCmdF(cmd *cobra.Command, _ []string) {
	epic, _ := cmd.Flags().GetString("epic")
	if err := campaign.Save(&model.Campaign{Epic: epic}); err != nil {
		ErrorAndExit(cmd, err)
	}
}
