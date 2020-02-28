package cmd

import (
	"git.ctrlz.es/mgdelacroix/campaigner/config"
	
	"github.com/spf13/cobra"
)

func TokenSetJiraCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "jira",
		Short: "Sets the value of the jira token",
		Args: cobra.ExactArgs(1),
		RunE: tokenSetJiraCmdF,
	}
}

func TokenSetGithubCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "github",
		Short: "Sets the value of the github token",
		Args: cobra.ExactArgs(1),
		RunE: tokenSetGithubCmdF,
	}
}

func TokenSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Sets the value of the platform tokens",
	}

	cmd.AddCommand(
		TokenSetJiraCmd(),
		TokenSetGithubCmd(),
	)

	return cmd
}

func TokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Subcommands related to tokens",
	}

	cmd.AddCommand(
		TokenSetCmd(),
	)

	return cmd
}

func tokenSetJiraCmdF(cmd *cobra.Command, args []string) error {
	cfg, err := config.ReadConfig()
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	cfg.JiraToken = args[0]
	if err := config.SaveConfig(cfg); err != nil {
		ErrorAndExit(cmd, err)
	}
	return nil
}

func tokenSetGithubCmdF(cmd *cobra.Command, args []string) error {
	cfg, err := config.ReadConfig()
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	cfg.GithubToken = args[0]
	if err := config.SaveConfig(cfg); err != nil {
		ErrorAndExit(cmd, err)
	}
	return nil
}
