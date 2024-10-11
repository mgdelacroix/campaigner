package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mgdelacroix/campaigner/app"
)

func StatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Prints the campaign status",
		Long:  "Prints the current status of the campaign and its tickets",
		Args:  cobra.NoArgs,
		RunE:  withAppE(statusCmdF),
	}
	cmd.Flags().Bool("md", false, "print output in markdown format")

	return cmd
}

func statusCmdF(a *app.App, cmd *cobra.Command, _ []string) error {
	md, err := cmd.Flags().GetBool("md")
	if err != nil {
		return err
	}

	a.Campaign.PrintStatus(md)

	return nil
}
