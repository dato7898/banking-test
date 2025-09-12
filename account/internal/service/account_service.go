package service

import (
	"account/internal/model"
	"account/internal/repository"
)

type AccountService interface {
	GetAccount(iban string) (*model.Account, error)
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
