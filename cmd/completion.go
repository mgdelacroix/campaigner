package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func CompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generates autocompletions for bash and zsh",
	}

	cmd.AddCommand(
		BashCompletionCmd(),
		ZshCompletionCmd(),
	)

	return cmd
}

func BashCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bash",
		Short: "Generates autocompletions for bash",
		Long: `Generates autocompletions for bash. To load them, run:

. <(campaigner completion bash)

To configure your bash shell to load completions for each session, add the above line to your ~/.bashrc`,
		Run: bashCompletionCmdF,
	}
}

func ZshCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "zsh",
		Short: "Generates autocompletions for zsh",
		Long: `Generates autocompletions for zsh. To load them, run:

. <(campaigner completion zsh)

To configure your zsh shell to load completions for each session, add the above line to your ~/.zshrc`,
		Run: zshCompletionCmdF,
	}
}

func bashCompletionCmdF(cmd *cobra.Command, args []string) {
	if err := cmd.Root().GenBashCompletion(cmd.OutOrStdout()); err != nil {
		ErrorAndExit(cmd, fmt.Errorf("unable to generate completions: %w", err))
	}
}

func zshCompletionCmdF(cmd *cobra.Command, args []string) {
	if err := cmd.Root().GenZshCompletion(cmd.OutOrStdout()); err != nil {
		ErrorAndExit(cmd, fmt.Errorf("unable to generate completions: %w", err))
	}
}
