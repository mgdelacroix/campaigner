package model

import (
	"fmt"
)

type Ticket map[string]interface{}

func RemoveDuplicateTickets(tickets []*Ticket, fileOnly bool) []*Ticket {
	ticketMap := map[string]*Ticket{}
	for _, t := range tickets {
		filename, _ := (*t)["filename"].(string)
		lineNo, _ := (*t)["lineNo"].(int)
		if fileOnly {
			ticketMap[filename] = t
		} else {
			ticketMap[fmt.Sprintf("%s:%d", filename, lineNo)] = t
		}
	}

	cleanTickets := []*Ticket{}
	for _, t := range ticketMap {
		cleanTickets = append(cleanTickets, t)
	}

	return cleanTickets
}
