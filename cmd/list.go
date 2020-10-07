package cmd

import (
	"github.com/spf13/cobra"

	"git.ctrlz.es/mgdelacroix/campaigner/app"
)

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Prints a list of the campaign's tickets",
		Long:  "Prints a list of the campaign's tickets with their statuses and external ids",
		Args:  cobra.NoArgs,
		Run:   withApp(listCmdF),
	}
}

func listCmdF(a *app.App, cmd *cobra.Command, _ []string) {
	a.Campaign.PrintList()
}
