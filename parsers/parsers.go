package parsers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"git.ctrlz.es/mgdelacroix/campaigner/model"
)

const (
	GREP  = "grep"
	GOVET = "govet"
)

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

func parseGovetLine(line string) (*model.Ticket, error) {
	parts := strings.Split(line, ":")
	if len(parts) < 4 {
		return nil, fmt.Errorf("cannot parse line: %s", line)
	}

	filename := parts[0]
	lineNo, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	text := strings.Join(parts[3:], "")

	return &model.Ticket{
		Data: map[string]interface{}{
			"filename": filename,
			"lineNo":   lineNo,
			"text":     strings.TrimSpace(text),
		},
	}, nil
}

func ParseWith(parser string) []*model.Ticket {
	var parseFn func(string) (*model.Ticket, error)
	switch parser {
	case GREP:
		parseFn = parseGrepLine
	case GOVET:
		parseFn = parseGovetLine
	default:
		fmt.Fprintf(os.Stderr, "Unknown parser %s", parser)
		return nil
	}

	tickets := []*model.Ticket{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ticket, _ := parseFn(scanner.Text())
		if ticket != nil {
			tickets = append(tickets, ticket)
		}
	}
	return tickets
}
