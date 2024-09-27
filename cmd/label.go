package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/mgdelacroix/campaigner/app"
)

const defaultEditor = "vim"

func ListLabelCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List the local campaign labels",
		Args:  cobra.NoArgs,
		Run:   withApp(listLabelCmdF),
	}
}

func RemoteLabelCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remote",
		Short: "List all the GitHub repository labels",
		Args:  cobra.NoArgs,
		Run:   withApp(remoteLabelCmdF),
	}
}

func UpdateLabelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates the campaign's GitHub labels",
		Args:  cobra.NoArgs,
		Run:   withApp(updateLabelCmdF),
	}

	cmd.Flags().BoolP("skip-check", "s", false, "do not check if the labels exist in the remote repository")

	return cmd
}

func LabelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Commands to manage GitHub labels",
	}

	cmd.AddCommand(
		ListLabelCmd(),
		RemoteLabelCmd(),
		UpdateLabelCmd(),
	)

	return cmd
}

func listLabelCmdF(a *app.App, cmd *cobra.Command, _ []string) {
	for _, label := range a.Campaign.Github.Labels {
		fmt.Println(label)
	}
}

func remoteLabelCmdF(a *app.App, cmd *cobra.Command, _ []string) {
	labels, err := a.ListLabels()
	if err != nil {
		ErrorAndExit(cmd, fmt.Errorf("cannot retrieve labels list: %w", err))
	}

	fmt.Printf("Labels for repository %s:\n\n", color.GreenString(a.Campaign.Github.Repo))
	for _, label := range labels {
		fmt.Println(label)
	}
}

func updateLabelCmdF(a *app.App, cmd *cobra.Command, _ []string) {
	skipCheck, _ := cmd.Flags().GetBool("skip-check")

	file, err := os.CreateTemp("", "campaigner-")
	if err != nil {
		ErrorAndExit(cmd, fmt.Errorf("cannot create temp file: %w", err))
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()
	labelBytes := []byte(strings.Join(a.Campaign.Github.Labels, "\n"))
	if _, err := file.Write(labelBytes); err != nil {
		ErrorAndExit(cmd, fmt.Errorf("cannot write labels to temp file: %w", err))
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = defaultEditor
	}

	editorCmd := exec.Command(editor, file.Name())
	editorCmd.Stdout = cmd.OutOrStdout()
	editorCmd.Stdin = cmd.InOrStdin()
	editorCmd.Stderr = cmd.ErrOrStderr()

	if err := editorCmd.Run(); err != nil {
		ErrorAndExit(cmd, fmt.Errorf("cannot run editor command: %w", err))
	}

	newLabelBytes, err := os.ReadFile(file.Name())
	if err != nil {
		ErrorAndExit(cmd, fmt.Errorf("cannot read file: %w", err))
	}

	newLabels := []string{}
	for _, label := range strings.Split(string(newLabelBytes), "\n") {
		if label != "" {
			newLabels = append(newLabels, strings.TrimSpace(label))
		}
	}

	if !skipCheck {
		ok, badLabels, err := a.CheckLabels(newLabels)
		if err != nil {
			ErrorAndExit(cmd, fmt.Errorf("cannot check new labels list: %w", err))
		}

		if !ok {
			ErrorAndExit(cmd, fmt.Errorf("these labels doesn't exist in the repository:\n\n%s", strings.Join(badLabels, "\n")))
		}
	}

	a.Campaign.Github.Labels = newLabels
	if err := a.Save(); err != nil {
		ErrorAndExit(cmd, fmt.Errorf("cannot save campaign: %w", err))
	}

	cmd.Println("Labels successfully updated")
}
