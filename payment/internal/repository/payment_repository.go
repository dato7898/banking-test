package repository

import (
	"payment/internal/model"

	"github.com/jmoiron/sqlx"
)

type PaymentRepository interface {
	Save(payment *model.Payment) error
	BeginTx() (*sqlx.Tx, error)
	WithTx(tx *sqlx.Tx) PaymentRepository
}

type paymentRepository struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewPaymentRepository(db *sqlx.DB) PaymentRepository {
	db.MustExec(`CREATE TABLE IF NOT EXISTS payments (
		id TEXT PRIMARY KEY,
		amount NUMERIC,
		from_account TEXT,
		to_account TEXT,
		user_id INT
	)`)
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Save(payment *model.Payment) error {
	query := `
		INSERT INTO payments (id, amount, from_account, to_account, user_id) 
		VALUES (:id, :amount, :from_account, :to_account, :user_id)
	`

	if r.tx != nil {
		_, err := r.tx.NamedExec(query, payment)
		return err
	}

	_, err := r.db.NamedExec(query, payment)
	return err
}

func (r *paymentRepository) BeginTx() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *paymentRepository) WithTx(tx *sqlx.Tx) PaymentRepository {
	return &paymentRepository{db: r.db, tx: tx}
}
