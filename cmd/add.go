package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"git.ctrlz.es/mgdelacroix/campaigner/campaign"
	"git.ctrlz.es/mgdelacroix/campaigner/model"
)

func GrepAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grep",
		Short: "Generates the tickets reading grep's output from stdin",
		Long: `Generates tickets for the campaign reading from the standard input the output grep. The grep command must be run with the -n flag. The generated ticket will contain three fields:

 - filename: the filename yield by grep
 - lineNo: the line number yield by grep
 - text: the trimmed line that grep captured for the expression
`,
		Example: `  grep -nriIF --include \*.go cobra.Command | campaigner add grep`,
		Args:    cobra.NoArgs,
		Run:     grepAddCmdF,
	}

	cmd.Flags().BoolP("file-only", "f", false, "Generates one ticket per file instead of per match")

	return cmd
}

func AgAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ag",
		Short:   "Generates the tickets reading ag's output from stdin",
		Long:    "Generates tickets for the campaign reading from the standard input the output ag",
		Example: `  ag cobra.Command | campaigner add ag`,
		Args:    cobra.NoArgs,
		RunE:    agAddCmdF,
	}

	cmd.Flags().BoolP("file-only", "f", false, "Generates one ticket per file instead of per match")

	return cmd
}

func GovetAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "govet",
		Short:   "Generates the tickets reading govet's output from stdin",
		Long:    "Generates tickets for the campaign reading from the standard input the output grep. The grep command must be run with the -json flag",
		Example: `  govet -json ./... | campaigner add govet`,
		Args:    cobra.NoArgs,
		RunE:    govetAddCmdF,
	}

	cmd.Flags().BoolP("file-only", "f", false, "Generates one ticket per file instead of per match")

	return cmd
}

func CsvAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "csv",
		Short:   "Generates the tickets reading a csv file",
		Example: `  campaigner add csv tickets.csv`,
		Args:    cobra.ExactArgs(1),
		Run:     csvAddCmdF,
	}
}

func AddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Adds tickets to the campaign",
		Long:  "Adds tickets to the campaign from the output of grep/ag/govet",
	}

	cmd.AddCommand(
		GrepAddCmd(),
		AgAddCmd(),
		GovetAddCmd(),
		CsvAddCmd(),
	)

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
		Data: map[string]interface{}{
			"filename": filename,
			"lineNo":   lineNo,
			"text":     strings.TrimSpace(text),
		},
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

func grepAddCmdF(cmd *cobra.Command, _ []string) {
	fileOnly, _ := cmd.Flags().GetBool("file-only")

	tickets := parseGrep()

	cmp, err := campaign.Read()
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	cmp.Tickets = append(cmp.Tickets, tickets...)
	cmp.Tickets = model.RemoveDuplicateTickets(cmp.Tickets, fileOnly)

	if err := campaign.Save(cmp); err != nil {
		ErrorAndExit(cmd, err)
	}
	cmd.Printf("%d tickets have been added\n", len(tickets))
}

func agAddCmdF(_ *cobra.Command, _ []string) error {
	return fmt.Errorf("not implemented yet")
}

func govetAddCmdF(_ *cobra.Command, _ []string) error {
	return fmt.Errorf("not implemented yet")
}

func csvAddCmdF(cmd *cobra.Command, args []string) {
	file, err := os.Open(args[0])
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	cmp, err := campaign.Read()
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	csvReader := csv.NewReader(bufio.NewReader(file))
	records, err := csvReader.ReadAll()
	if err != nil {
		ErrorAndExit(cmd, err)
	}

	headers := records[0]
	for _, line := range records[1:] {
		data := map[string]interface{}{}
		for i, header := range headers {
			data[header] = line[i]
		}
		cmp.Tickets = append(cmp.Tickets, &model.Ticket{Data: data})
	}

	if err := campaign.Save(cmp); err != nil {
		ErrorAndExit(cmd, err)
	}
	cmd.Printf("%d tickets have been added\n", len(records[1:]))
}
