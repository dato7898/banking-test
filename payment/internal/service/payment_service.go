package service

import (
	"payment/internal/model"
	"payment/internal/repository"
	"payment/internal/validator"

	"github.com/google/uuid"

	pb "payment/proto"
)

type PaymentService interface {
	CreatePayment(userID int, req model.CreatePaymentRequest) (string, error)
}

type paymentService struct {
	paymentRepo      repository.PaymentRepository
	paymentValidator validator.PaymentValidator
	accountClient    repository.AccountClient
	producer         repository.KafkaProducer
}

func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	paymentValidator validator.PaymentValidator,
	accountClient repository.AccountClient,
	producer repository.KafkaProducer,
) PaymentService {
	return &paymentService{
		paymentRepo:      paymentRepo,
		paymentValidator: paymentValidator,
		accountClient:    accountClient,
		producer:         producer,
	}
}

func (s *paymentService) CreatePayment(userID int, req model.CreatePaymentRequest) (string, error) {
	id := uuid.New().String()

	if err := s.paymentValidator.Validate(req, userID); err != nil {
		return "", err
	}

	_, err := s.accountClient.Withdrawal(&pb.OperationRequest{Iban: req.From, Amount: req.Amount})
	if err != nil {
		return "", err
	}

	if err = s.processPaymentCreation(userID, id, req); err != nil {
		// TODO переделать на вызов топика с ретраем
		_, accErr := s.accountClient.Replenishment(&pb.OperationRequest{Iban: req.From, Amount: req.Amount})
		if accErr != nil {
			return "", accErr
		}
		return "", err
	}

	return id, nil
}

func (s *paymentService) processPaymentCreation(userID int, id string, req model.CreatePaymentRequest) error {
	tx, err := s.paymentRepo.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	repo := s.paymentRepo.WithTx(tx)

	payment := &model.Payment{
		ID:     id,
		Amount: req.Amount,
		From:   req.From,
		To:     req.To,
		UserID: userID,
	}

	if err := repo.Save(payment); err != nil {
		return err
	}

	if err := s.producer.Publish(nil, map[string]string{"id": id}); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
