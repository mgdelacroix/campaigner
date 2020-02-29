package cmd

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"git.ctrlz.es/mgdelacroix/campaigner/model"
)

const defaultGrepOpts = "-nrFI"

func AddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Adds tickets to the campaign",
		Args:  cobra.NoArgs,
		Run:   addCmdF,
	}

	cmd.Flags().StringP("dir", "d", "", "directory containing the source code")
	_ = cmd.MarkFlagRequired("dir")
	cmd.Flags().StringSliceP("grep", "g", []string{}, "runs a grep command to generate the tickets")
	cmd.Flags().BoolP("case-insensitive", "i", false, "makes the search case insensitive")
	// cmd.Flags().StringP("govet", "v", "", "runs a govet command to generate the tickets")
	// govet bin path?

	return cmd
}

func parseLine(line string) (*model.Ticket, error) {
	parts := strings.Split(line, ":")
	if len(parts) < 3 {
		return nil, fmt.Errorf("cannot parse line: %s", line)
	}

	filename := parts[0]
	lineNo, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	text := strings.Join(parts[2:], "")

	return &model.Ticket{filename, lineNo, text}, nil
}

func RunGrep(dir, str string, caseInsensitive bool) ([]*model.Ticket, error) {
	opts := defaultGrepOpts
	if caseInsensitive {
		opts = opts + "i"
	}

	out, err := exec.Command("grep", opts, str, dir).Output()
	if err != nil {
		return nil, err
	}

	tickets := []*model.Ticket{}
	for _, line := range strings.Split(string(out), "\n") {
		// ToDo: get and check error
		ticket, _ := parseLine(line)
		if ticket != nil {
			tickets = append(tickets, ticket)
		}
	}

	return tickets, nil
}

func RunGreps(dir string, strs []string, caseInsensitive bool) ([]*model.Ticket, error) {
	tickets := []*model.Ticket{}
	for _, str := range strs {
		results, err := RunGrep(dir, str, caseInsensitive)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, results...)
	}

	return tickets, nil
}

func addCmdF(cmd *cobra.Command, _ []string) {
	dir, _ := cmd.Flags().GetString("dir")
	grepStrs, _ := cmd.Flags().GetStringSlice("grep")
	caseInsensitive, _ := cmd.Flags().GetBool("case-insensitive")

	tickets, err := RunGreps(dir, grepStrs, caseInsensitive)
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	for _, ticket := range tickets {
		fmt.Printf("%+v\n", ticket)
	}
}
