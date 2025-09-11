package service

import (
	pb "payment-worker/proto"

	"encoding/json"
	"log"
	"payment-worker/internal/repository"
)

type PaymentProcess interface {
	ProcessMessage(msg []byte) error
}

type paymentProcess struct {
	repo   repository.PaymentRepository
	client repository.PaymentIntegrationClient
}

func NewPaymentProcessor(repo repository.PaymentRepository, client repository.PaymentIntegrationClient) PaymentProcess {
	return &paymentProcess{repo: repo, client: client}
}

func (p *paymentProcess) ProcessMessage(msg []byte) error {
	var payload map[string]string
	if err := json.Unmarshal(msg, &payload); err != nil {
		log.Println("Invalid Kafka payload", err)
		return err
	}

	id := payload["id"]
	payment, err := p.repo.FindByID(id)
	if err != nil {
		log.Println("Payment not found", id, err)
		return err
	}

	log.Println("📩 Processing payment", payment)

	return p.client.ProcessPayment(&pb.Payment{
		Id:     payment.ID,
		Amount: payment.Amount,
		To:     payment.To,
	})
}
