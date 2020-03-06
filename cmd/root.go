package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "campaigner",
		Short: "Create and manage Open Source campaigns",
	}

	cmd.AddCommand(
		AddCmd(),
		FilterCmd(),
		InitCmd(),
		StandaloneCmd(),
		StatusCmd(),
		TokenCmd(),
		PublishCmd(),
		SyncCmd(),
	)

	return cmd
}

func Execute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
