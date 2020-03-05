package model

type Campaign struct {
	Url     string    `json:"url"`
	Project string    `json:"project"`
	Epic    string    `json:"epic"`
	Summary string    `json:"summary"`
	Tickets []*Ticket `json:"tickets,omitempty"`
}
