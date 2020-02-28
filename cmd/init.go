package cmd

import (
	"github.com/spf13/cobra"
)

func InitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Creates a new campaign in the current directory",
		Args:  cobra.NoArgs,
		RunE:  initCmdF,
	}

	// add mandatory flags for epic, tags, etc
}

func initCmdF(_ *cobra.Command, _ []string) error {
	// creates the campaign.json file
	return nil
}
