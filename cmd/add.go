package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mgdelacroix/campaigner/app"
	"github.com/mgdelacroix/campaigner/model"
	"github.com/mgdelacroix/campaigner/parsers"
)

func GrepAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grep",
		Short: "Generates the tickets reading grep's output from stdin",
		Long: `Generates tickets for the campaign reading the output of grep from the standard input. The grep command must be run with the -n flag. The generated ticket will contain three fields:

 - filename: the filename yield by grep
 - lineNo: the line number yield by grep
 - text: the trimmed line that grep captured for the expression
`,
		Example: `  grep -nriIF --include \*.go cobra.Command | campaigner add grep`,
		Args:    cobra.NoArgs,
		Run:     withApp(grepAddCmdF),
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
		Use:   "govet",
		Short: "Generates the tickets reading govet's output from stdin",
		Long: `Generates tickets for the campaign reading the output of govet from the standard input. Govet usually writes to the standard error descriptor, so the output must be redirected. The generated ticket will contain three fields:

 - filename: the filename yield by grep
 - lineNo: the line number yield by grep
 - text: the text containing the govet error
`,
		Example: `  govet ./... 2>&1 | campaigner add govet`,
		Args:    cobra.NoArgs,
		Run:     withApp(govetAddCmdF),
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
		Run:     withApp(csvAddCmdF),
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

func grepAddCmdF(a *app.App, cmd *cobra.Command, _ []string) {
	fileOnly, _ := cmd.Flags().GetBool("file-only")

	tickets := parsers.ParseWith(parsers.GREP)
	addedTickets := a.Campaign.AddTickets(tickets, fileOnly)

	if err := a.Save(); err != nil {
		ErrorAndExit(cmd, err)
	}
	cmd.Printf("%d tickets have been added\n", addedTickets)
}

func agAddCmdF(_ *cobra.Command, _ []string) error {
	return fmt.Errorf("not implemented yet")
}

func govetAddCmdF(a *app.App, cmd *cobra.Command, _ []string) {
	fileOnly, _ := cmd.Flags().GetBool("file-only")

	tickets := parsers.ParseWith(parsers.GOVET)
	addedTickets := a.Campaign.AddTickets(tickets, fileOnly)

	if err := a.Save(); err != nil {
		ErrorAndExit(cmd, err)
	}
	cmd.Printf("%d tickets have been added\n", addedTickets)
}

func csvAddCmdF(a *app.App, cmd *cobra.Command, args []string) {
	file, err := os.Open(args[0])
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
		a.Campaign.Tickets = append(a.Campaign.Tickets, &model.Ticket{Data: data})
	}

	if err := a.Save(); err != nil {
		ErrorAndExit(cmd, err)
	}
	cmd.Printf("%d tickets have been added\n", len(records[1:]))
}
