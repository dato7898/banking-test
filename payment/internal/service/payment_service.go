package service

import (
	"errors"
	"payment/internal/model"
	"payment/internal/repository"

	"github.com/google/uuid"

	pb "payment/proto"
)

type PaymentService interface {
	CreatePayment(userID int, req model.CreatePaymentRequest) (string, error)
}

type paymentService struct {
	paymentRepo   repository.PaymentRepository
	accountClient repository.AccountClient
	producer      repository.KafkaProducer
}

func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	accountClient repository.AccountClient,
	producer repository.KafkaProducer,
) PaymentService {
	return &paymentService{paymentRepo: paymentRepo, accountClient: accountClient, producer: producer}
}

func (s *paymentService) CreatePayment(userID int, req model.CreatePaymentRequest) (string, error) {
	id := uuid.New().String()

	accountReq := &pb.GetAccountRequest{
		Iban: req.From,
	}

	account, err := s.accountClient.GetAccount(accountReq)
	if err != nil {
		return "", err
	}

	if account.UserID != int32(userID) {
		return "", errors.New("account not found")
	}

	if account.Amount < req.Amount {
		return "", errors.New("insufficient funds")
	}

	accountReq.Iban = req.To

	account, err = s.accountClient.GetAccount(accountReq)
	if err != nil {
		return "", err
	}

	payment := &model.Payment{
		ID:     id,
		Amount: req.Amount,
		From:   req.From,
		To:     req.To,
		UserID: userID,
	}

	if err := s.paymentRepo.Save(payment); err != nil {
		return "", err
	}

	if err := s.producer.Publish("payments", nil, map[string]string{"id": id}); err != nil {
		return "", err
	}

	return id, nil
}
