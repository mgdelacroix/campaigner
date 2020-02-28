package cmd

import (
	"github.com/spf13/cobra"
)

func AddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Adds tickets to the campaign",
		Args: cobra.NoArgs,
		RunE: addCmdF,
	}

	// add flags and examples
}

func addCmdF(_ *cobra.Command, _ []string) error {
	return nil
}
