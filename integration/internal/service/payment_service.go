package service

import (
	"integration/internal/model"
	"integration/internal/repository"
	"log"
)

type PaymentService interface {
	ProcessPayment(p *model.Payment) error
}

type paymentService struct {
	external repository.ExternalSystem
}

func NewPaymentService(external repository.ExternalSystem) PaymentService {
	return &paymentService{external: external}
}

func (s *paymentService) ProcessPayment(p *model.Payment) error {
	log.Println("➡️ Sending payment to external system:", p)
	return s.external.SendPayment(p)
}
