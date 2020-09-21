package cmd

import (
	"git.ctrlz.es/mgdelacroix/campaigner/app"

	"github.com/spf13/cobra"
)

func SyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Syncs the tickets",
		Long:  "Synchronizes the status of the published tickets with remote providers",
		Args:  cobra.NoArgs,
		Run:   withApp(syncCmdF),
	}
}

func syncCmdF(a *app.App, cmd *cobra.Command, _ []string) {
	if err := a.GithubSync(); err != nil {
		ErrorAndExit(cmd, err)
	}
	cmd.Println("Synchronization completed")
}
