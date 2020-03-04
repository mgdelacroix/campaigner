package cmd

import (
	"github.com/spf13/cobra"
)

func SyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronizes the status of the tickets with remote providers",
		Args:  cobra.NoArgs,
		Run:   syncCmdF,
	}

	return cmd
}

func syncCmdF(_ *cobra.Command, _ []string) {}
