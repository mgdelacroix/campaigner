package model

type Ticket struct {
	Filename string `json:"filename"`
	Line     int    `json:"line"`
	Text     string `json:"text"` // needed??
}
