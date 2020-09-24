package model

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func (c *Campaign) PrintUserReport() {
	userTickets := map[string]int{}
	for _, ticket := range c.Tickets {
		user := ticket.GithubAssignee
		if user != "" {
			if count, ok := userTickets[user]; ok {
				userTickets[user] = count+1
			} else {
				userTickets[user] = 1
			}
		}
	}

	fmt.Println("User report:\n")
	w := tabwriter.NewWriter(os.Stdout, 3, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "GitHub username\tTickets closed\t")
	for user, count := range userTickets {
		fmt.Fprintf(w, "%s\t%d\t\n", user, count)
	}
	w.Flush()
	fmt.Println()
}
