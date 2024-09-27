package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mgdelacroix/campaigner/app"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Prints a list of the campaign's tickets",
		Long:  "Prints a list of the campaign's tickets with their statuses and external ids",
		Args:  cobra.NoArgs,
		Run:   withApp(listCmdF),
	}

	cmd.Flags().BoolP("published-only", "p", false, "list only published tickets")
	cmd.Flags().BoolP("links", "l", false, "print full links for jira and github ticket numbers")

	return cmd
}

func listCmdF(a *app.App, cmd *cobra.Command, _ []string) {
	publishedOnly, _ := cmd.Flags().GetBool("published-only")
	printLinks, _ := cmd.Flags().GetBool("links")

	a.Campaign.PrintList(publishedOnly, printLinks)
}
