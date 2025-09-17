package repository

import (
	"account/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AccountRepository interface {
	GetAccount(iban string) (*model.Account, error)
	GetForUpdate(iban string) (*model.Account, error)
	Replenishment(iban string, amount float64) error
	Withdrawal(iban string, amount float64) error
	BeginTx() (*sqlx.Tx, error)
	WithTx(tx *sqlx.Tx) AccountRepository
}

type accountRepository struct {
	db *sqlx.DB
	tx *sqlx.Tx
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

func (r *accountRepository) GetForUpdate(iban string) (*model.Account, error) {
	var a model.Account
	query := "SELECT id, iban, amount, user_id FROM account WHERE iban=$1 FOR NO KEY UPDATE"

	if r.tx != nil {
		err := r.tx.Get(&a, query, iban)
		if err != nil {
			return nil, err
		}
		return &a, nil
	}

	err := r.db.Get(&a, query, iban)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *accountRepository) Replenishment(iban string, amount float64) error {
	_, err := r.db.NamedExec(`UPDATE account SET amount = (amount + :amount) WHERE iban=:iban`, &model.OperationRequest{
		Iban: iban, Amount: amount,
	})
	return err
}

func (r *accountRepository) Withdrawal(iban string, amount float64) error {
	query := `UPDATE account SET amount = (amount - :amount) WHERE iban=:iban`
	params := &model.OperationRequest{Iban: iban, Amount: amount}

	if r.tx != nil {
		_, err := r.tx.NamedExec(query, params)
		return err
	}

	_, err := r.db.NamedExec(query, params)
	return err
}

func (r *accountRepository) BeginTx() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *accountRepository) WithTx(tx *sqlx.Tx) AccountRepository {
	return &accountRepository{db: r.db, tx: tx}
}
