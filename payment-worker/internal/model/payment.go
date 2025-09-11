package model

type Payment struct {
	ID     string  `db:"id"`
	Amount float64 `db:"amount"`
	To     string  `db:"to_account"`
}
