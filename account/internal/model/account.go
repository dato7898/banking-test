package model

type Account struct {
	ID     string  `db:"id"`
	Iban   string  `db:"iban"`
	Amount float64 `db:"amount"`
	UserID int     `db:"user_id"`
}
