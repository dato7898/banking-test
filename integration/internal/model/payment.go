package model

type Payment struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
	To     string  `json:"to"`
}
