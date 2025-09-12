package model

type Payment struct {
	ID     string  `db:"id" json:"id"`
	Amount float64 `db:"amount" json:"amount"`
	From   string  `db:"from_account" json:"from"`
	To     string  `db:"to_account" json:"to"`
	UserID int     `db:"user_id" json:"userID"`
}

type CreatePaymentRequest struct {
	Amount float64 `json:"amount"`
	From   string  `json:"from"`
	To     string  `json:"to"`
}
