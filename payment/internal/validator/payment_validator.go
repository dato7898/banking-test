package validator

import (
	"errors"
	"payment/internal/model"
	"payment/internal/repository"

	pb "payment/proto"
)

type PaymentValidator interface {
	Validate(req model.CreatePaymentRequest, userID int) error
}

type paymentValidator struct {
	accountClient repository.AccountClient
}

func NewPaymentValidator(accountClient repository.AccountClient) PaymentValidator {
	return &paymentValidator{accountClient: accountClient}
}

func (v *paymentValidator) Validate(req model.CreatePaymentRequest, userID int) error {
	fromAcc, err := v.accountClient.GetAccount(&pb.GetAccountRequest{Iban: req.From})
	if err != nil {
		return errors.New("debit account not found")
	}

	if fromAcc.UserID != int32(userID) {
		return errors.New("debit account not found")
	}

	if fromAcc.Amount < req.Amount {
		return errors.New("insufficient funds")
	}

	_, err = v.accountClient.GetAccount(&pb.GetAccountRequest{Iban: req.To})
	if err != nil {
		return errors.New("destination account not found")
	}

	return nil
}
