package repository

import (
	"payment-worker/internal/model"

	"github.com/jmoiron/sqlx"
)

type PaymentRepository interface {
	FindByID(id string) (*model.Payment, error)
}

type paymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) FindByID(id string) (*model.Payment, error) {
	var p model.Payment
	err := r.db.Get(&p, "SELECT id, amount, to_account FROM payments WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
