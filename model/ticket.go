package model

import (
	"fmt"
)

type Ticket struct {
	Filename string `json:"filename"`
	LineNo   int    `json:"line_no"`
	Text     string `json:"text"`
}

func RemoveDuplicateTickets(tickets []*Ticket, fileOnly bool) []*Ticket {
	ticketMap := map[string]*Ticket{}
	for _, t := range tickets {
		if fileOnly {
			t.Text = ""
			t.LineNo = 0
			ticketMap[t.Filename] = t
		} else {
			ticketMap[fmt.Sprintf("%s:%d", t.Filename, t.LineNo)] = t
		}
	}

	cleanTickets := []*Ticket{}
	for _, t := range ticketMap {
		cleanTickets = append(cleanTickets, t)
	}

	return cleanTickets
}
