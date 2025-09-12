package service

import (
	"payment/internal/model"
	"payment/internal/repository"

	"github.com/google/uuid"
)

type PaymentService interface {
	CreatePayment(userID int, req model.CreatePaymentRequest) (string, error)
}

type paymentService struct {
	repo     repository.PaymentRepository
	producer repository.KafkaProducer
}

func NewPaymentService(repo repository.PaymentRepository, producer repository.KafkaProducer) PaymentService {
	return &paymentService{repo: repo, producer: producer}
}

func (s *paymentService) CreatePayment(userID int, req model.CreatePaymentRequest) (string, error) {
	id := uuid.New().String()
	payment := &model.Payment{
		ID:     id,
		Amount: req.Amount,
		From:   req.From,
		To:     req.To,
		UserID: userID,
	}

	if err := s.repo.Save(payment); err != nil {
		return "", err
	}

	if err := s.producer.Publish("payments", nil, map[string]string{"id": id}); err != nil {
		return "", err
	}

	return id, nil
}
