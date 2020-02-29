package model

type Campaign struct {
	Epic    string    `json:"epic"`
	Summary string    `json:"summary"`
	Tickets []*Ticket `json:"tickets,omitempty"`
}
