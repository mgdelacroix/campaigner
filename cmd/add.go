package cmd

import (
	"github.com/spf13/cobra"

	"git.ctrlz.es/mgdelacroix/campaigner/model"
)

func AddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Adds tickets to the campaign",
		Args:  cobra.NoArgs,
		RunE:  addCmdF,
	}

	cmd.Flags().StringP("dir", "d", "", "directory containing the source code")
	_ = cmd.MarkFlagRequired("dir")
	cmd.Flags().StringSliceP("grep", "g", []string{}, "runs a grep command to generate the tickets")
	cmd.Flags().BoolP("case-insensitive", "i", false, "makes the search case insensitive")
	// cmd.Flags().StringP("govet", "v", "", "runs a govet command to generate the tickets")
	// govet bin path?

	return cmd
}

func RunGrep(dir string, strs []string, caseInsensitive bool) ([]model.Ticket, error) {
	// grep -nrI TEXT .
	// -i as well with case insensitive

	return nil, nil
}

func addCmdF(cmd *cobra.Command, _ []string) error {
	// either govet or grep

	return nil
}
