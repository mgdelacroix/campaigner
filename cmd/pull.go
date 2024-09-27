package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mgdelacroix/campaigner/app"
)

func PullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Imports tickets from Jira",
		Long:  "Imports all tickets from a Jira epic issue. This command is only intended to use when the Jira tickets are already created and we want to use campaigner to import and manage them",
		Args:  cobra.NoArgs,
		Run:   withApp(pullCmdF),
	}
}

func pullCmdF(a *app.App, cmd *cobra.Command, _ []string) {
	tickets, err := a.GetTicketsFromJiraEpic()
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	addedTickets := a.Campaign.AddTickets(tickets, false)

	if err := a.Save(); err != nil {
		ErrorAndExit(cmd, err)
	}
	cmd.Printf("%d tickets have been added\n", addedTickets)
}
