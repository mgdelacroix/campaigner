package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func ErrorAndExit(cmd *cobra.Command, err error) {
	cmd.PrintErrln("ERROR: " + err.Error())
	os.Exit(1)
}
