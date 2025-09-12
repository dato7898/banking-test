package repository

import (
	"account/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AccountRepository interface {
	GetAccount(iban string) (*model.Account, error)
}

type accountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) AccountRepository {
	db.MustExec(`CREATE TABLE IF NOT EXISTS account (
		id TEXT PRIMARY KEY,
		iban TEXT NOT NULL UNIQUE,
		amount NUMERIC NOT NULL,
		user_id INT NOT NULL
	)`)

	db.NamedExec(`INSERT INTO account (id, iban, amount, user_id) VALUES (:id, :iban, :amount, :user_id) ON CONFLICT DO NOTHING`, &model.Account{
		ID:     uuid.New().String(),
		Iban:   "12345",
		Amount: 100,
		UserID: 1,
	})

	db.NamedExec(`INSERT INTO account (id, iban, amount, user_id) VALUES (:id, :iban, :amount, :user_id) ON CONFLICT DO NOTHING`, &model.Account{
		ID:     uuid.New().String(),
		Iban:   "23456",
		Amount: 100,
		UserID: 2,
	})

	return &accountRepository{db: db}
}

func (r *accountRepository) GetAccount(iban string) (*model.Account, error) {
	var a model.Account
	err := r.db.Get(&a, "SELECT id, iban, amount, user_id FROM account WHERE iban=$1", iban)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
