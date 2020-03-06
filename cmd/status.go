package cmd

import (
	"git.ctrlz.es/mgdelacroix/campaigner/campaign"

	"github.com/spf13/cobra"
)

func StatusCmd() *cobra.Command {
	return &cobra.Command{
		Use: "status",
		Short: "Prints the current status of the campaign",
		Args: cobra.NoArgs,
		Run: statusCmdF,
	}
}

func statusCmdF(cmd *cobra.Command, _ []string) {
	cmp, err := campaign.Read()
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	cmp.PrintStatus(cmd.OutOrStdout())
}
