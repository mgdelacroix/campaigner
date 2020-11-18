package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func createTicket(filename string, lineNo int) *Ticket {
	return &Ticket{
		Data: map[string]interface{}{
			"filename": filename,
			"lineNo":   lineNo,
		},
	}
}

func TestAddTickets(t *testing.T) {
	t.Run("Tickets should be added correctly with fileOnly disabled", func(t *testing.T) {
		campaign := &Campaign{}
		ticketsToAdd := []*Ticket{
			createTicket("user.txt", 1),
			createTicket("test.txt", 2),
			createTicket("sample.txt", 3),
			createTicket("user.txt", 4),
		}

		addedTickets := campaign.AddTickets(ticketsToAdd, false)
		require.Equal(t, 4, addedTickets)
		require.Len(t, campaign.Tickets, 4)
	})

	t.Run("Should account for already existing tickets with fileOnly disabled", func(t *testing.T) {
		campaign := &Campaign{
			Tickets: []*Ticket{
				createTicket("user.txt", 1),
				createTicket("test.txt", 2),
			},
		}
		ticketsToAdd := []*Ticket{
			createTicket("user.txt", 1),
			createTicket("test.txt", 2),
			createTicket("sample.txt", 3),
			createTicket("user.txt", 4),
		}

		addedTickets := campaign.AddTickets(ticketsToAdd, false)
		require.Equal(t, 2, addedTickets)
		require.Len(t, campaign.Tickets, 4)
	})

	t.Run("Tickets should be added correctly with fileOnly enabled", func(t *testing.T) {
		campaign := &Campaign{}
		ticketsToAdd := []*Ticket{
			createTicket("user.txt", 1),
			createTicket("test.txt", 2),
			createTicket("sample.txt", 3),
			createTicket("user.txt", 4),
		}

		addedTickets := campaign.AddTickets(ticketsToAdd, true)
		require.Equal(t, 3, addedTickets)
		require.Len(t, campaign.Tickets, 3)
	})

	t.Run("Should account for already existing tickets with fileOnly enabled", func(t *testing.T) {
		campaign := &Campaign{
			Tickets: []*Ticket{
				createTicket("user.txt", 1),
				createTicket("test.txt", 2),
			},
		}
		ticketsToAdd := []*Ticket{
			createTicket("user.txt", 1),
			createTicket("test.txt", 2),
			createTicket("sample.txt", 3),
			createTicket("user.txt", 4),
		}

		addedTickets := campaign.AddTickets(ticketsToAdd, true)
		require.Equal(t, 1, addedTickets)
		require.Len(t, campaign.Tickets, 3)
	})
}
