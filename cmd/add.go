package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"git.ctrlz.es/mgdelacroix/campaigner/campaign"
	"git.ctrlz.es/mgdelacroix/campaigner/model"
)

func AddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Adds tickets to the campaign from the output of grep/ag/govet",
		Long: `Generates tickets for the campaign reading from the standard input the output of one of the following three commands:
  - grep (should be run with the -n flag)
  - ag
  - govet (should be run with the -json flag)`,
		Example: `  grep -nriIF --include \*.go cobra.Command | campaigner add --grep
  ag cobra.Command | campaigner add --ag
  govet -json ./... | campaigner add --govet`,
		Args: cobra.NoArgs,
		RunE: addCmdF,
	}

	cmd.Flags().BoolP("ag", "a", false, "generates the tickets reading ag's output from stdin")
	cmd.Flags().BoolP("grep", "g", false, "generates the tickets reading grep's output from stdin")
	cmd.Flags().BoolP("govet", "v", false, "generates the tickets reading govet's output from stdin")

	return cmd
}

func parseGrepLine(line string) (*model.Ticket, error) {
	// ToDo: it would be great to be able to relate a line with its
	// parent method, at least for JS and Golang
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

	return &model.Ticket{
		Filename: filename,
		LineNo:   lineNo,
		Text:     text,
	}, nil
}

func parseGrep() []*model.Ticket {
	tickets := []*model.Ticket{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ticket, _ := parseGrepLine(scanner.Text())
		if ticket != nil {
			tickets = append(tickets, ticket)
		}
	}
	return tickets
}

func addCmdF(cmd *cobra.Command, _ []string) error {
	grep, _ := cmd.Flags().GetBool("grep")
	ag, _ := cmd.Flags().GetBool("ag")
	govet, _ := cmd.Flags().GetBool("govet")

	if !grep && !ag && !govet {
		return fmt.Errorf("one of --grep --ag --govet flags should be active")
	}

	var tickets []*model.Ticket
	switch {
	case grep:
		tickets = parseGrep()
	default:
		return fmt.Errorf("not implemented yet")
	}

	cmp, err := campaign.Read()
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	cmp.Tickets = append(cmp.Tickets, tickets...)
	cmp.Tickets = model.RemoveDuplicateTickets(cmp.Tickets)

	if err := campaign.Save(cmp); err != nil {
		ErrorAndExit(cmd, err)
	}
	return nil
}
