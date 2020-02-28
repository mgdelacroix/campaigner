package cmd

import (
	"github.com/spf13/cobra"
)

func TokenSetJiraCmd() *cobra.Command {
	return &cobra.Command{
		Use: "jira",
		Short: "Sets the value of the jira token",
	}
}

func TokenSetGithubCmd() *cobra.Command {
	return &cobra.Command{
		Use: "github",
		Short: "Sets the value of the github token",
	}
}

func TokenSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "set",
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
		Use: "token",
		Short: "Subcommands related to tokens",
	}

	cmd.AddCommand(
		TokenSetCmd(),
	)

	return cmd
}
