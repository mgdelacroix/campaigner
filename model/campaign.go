package model

type Campaign struct {
	Epic    string    `json:"epic"`
	Tickets []*Ticket `json:"tickets,omitempty"`
}
