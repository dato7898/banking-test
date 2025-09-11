package repository

import (
	"payment/internal/model"

	"github.com/jmoiron/sqlx"
)

type PaymentRepository interface {
	Save(payment *model.Payment) error
}

type paymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) PaymentRepository {
	db.MustExec(`CREATE TABLE IF NOT EXISTS payments (
		id TEXT PRIMARY KEY,
		amount NUMERIC,
		currency TEXT,
		from_account TEXT,
		to_account TEXT,
		user_id INT
	)`)
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Save(payment *model.Payment) error {
	_, err := r.db.NamedExec(`
		INSERT INTO payments (id, amount, currency, from_account, to_account, user_id) 
		VALUES (:id, :amount, :currency, :from_account, :to_account, :user_id)
	`, payment)
	return err
}
