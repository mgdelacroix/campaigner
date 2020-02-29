package model

type Ticket struct {
	Filename string `json:"filename"`
	LineNo   int    `json:"line_no"`
	Text     string `json:"text"` // needed??
}
