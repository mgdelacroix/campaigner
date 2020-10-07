package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"git.ctrlz.es/mgdelacroix/campaigner/app"
)

func withApp(f func(*app.App, *cobra.Command, []string)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		campaignPath, _ := cmd.Flags().GetString("campaign")

		a, err := app.NewApp(campaignPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
			os.Exit(1)
		}

		f(a, cmd, args)
	}
}

func withAppE(f func(*app.App, *cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		campaignPath, _ := cmd.Flags().GetString("campaign")

		a, err := app.NewApp(campaignPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR: "+err.Error())
			os.Exit(1)
		}

		return f(a, cmd, args)
	}
}

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "campaigner",
		Short: "Create and manage Open Source campaigns",
	}

	cmd.PersistentFlags().String("campaign", "campaign.json", "the path to the campaign file")

	cmd.AddCommand(
		AddCmd(),
		// FilterCmd(),
		InitCmd(),
		LabelCmd(),
		StatusCmd(),
		PublishCmd(),
		PullCmd(),
		SyncCmd(),
		ReportCmd(),
		CompletionCmd(),
	)

	return cmd
}

func Execute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
