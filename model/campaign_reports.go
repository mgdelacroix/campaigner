package model

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

type Contritutor struct {
	Username      string
	Contributions int
}

func (c *Campaign) PrintUserReport() {
	contributors := c.Contritutor()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "GitHub username\tTickets closed\t")
	fmt.Fprintln(w, "---------------\t--------------\t")

	for _, c := range contributors {
		fmt.Fprintf(w, "%s\t%d\t\n", c.Username, c.Contributions)
	}

	w.Flush()
	fmt.Println()
}

func (c *Campaign) Contritutor() []Contritutor {
	contributions := map[string]int{}
	for _, ticket := range c.Tickets {
		if !ticket.IsClosed() {
			continue
		}

		user := ticket.GithubAssignee
		if user != "" {
			if count, ok := contributions[user]; ok {
				contributions[user] = count + 1
			} else {
				contributions[user] = 1
			}
		}
	}

	contributors := make([]Contritutor, 0, len(contributions))
	for user, count := range contributions {
		contributors = append(contributors, Contritutor{Username: user, Contributions: count})
	}
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Contributions > contributors[j].Contributions
	})

	return contributors
}
