package model

// ToDo: add key-value extra params as a map to allow for customfield_whatever = team
type Campaign struct {
	Url       string    `json:"url"`
	Project   string    `json:"project"`
	Epic      string    `json:"epic"`
	IssueType string    `json:"issue_type"`
	Summary   string    `json:"summary"`
	Template  string    `json:"template"`
	Tickets   []*Ticket `json:"tickets,omitempty"`
}

func (c *Campaign) NextUnpublishedTicket() *Ticket {
	for _, ticket := range c.Tickets {
		if ticket.JiraLink == "" {
			return ticket
		}
	}
	return nil
}
