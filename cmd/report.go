package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mgdelacroix/campaigner/app"
)

func UsersReportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "users",
		Short: "A users report",
		Args:  cobra.NoArgs,
		Run:   withApp(userReportCmdF),
	}
}

func ReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generates reports on campaign information",
	}

	cmd.AddCommand(
		UsersReportCmd(),
	)

	return cmd
}

func userReportCmdF(a *app.App, cmd *cobra.Command, _ []string) {
	a.Campaign.PrintUserReport()
}
