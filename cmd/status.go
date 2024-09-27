package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mgdelacroix/campaigner/app"
)

func StatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Prints the campaign status",
		Long:  "Prints the current status of the campaign and its tickets",
		Args:  cobra.NoArgs,
		Run:   withApp(statusCmdF),
	}
}

func statusCmdF(a *app.App, cmd *cobra.Command, _ []string) {
	a.Campaign.PrintStatus()
}
