package service

import (
	"account/internal/model"
	"account/internal/repository"
	"errors"
)

type AccountService interface {
	GetAccount(iban string) (*model.Account, error)
	Replenishment(iban string, amount float64) error
	Withdrawal(iban string, amount float64) error
}

type accountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return &accountService{accountRepo: accountRepo}
}

func (s *accountService) GetAccount(iban string) (*model.Account, error) {
	return s.accountRepo.GetAccount(iban)
}

func (s *accountService) Replenishment(iban string, amount float64) error {
	return s.accountRepo.Replenishment(iban, amount)
}

func (s *accountService) Withdrawal(iban string, amount float64) error {
	tx, err := s.accountRepo.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	repo := s.accountRepo.WithTx(tx)

	acc, err := repo.GetForUpdate(iban)
	if err != nil {
		return err
	}

	if acc.Amount-amount < 0 {
		return errors.New("insufficient funds")
	}

	err = repo.Withdrawal(iban, amount)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
