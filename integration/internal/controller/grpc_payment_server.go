package controller

import (
	"context"
	"integration/internal/model"
	"integration/internal/service"
	pb "integration/proto"
)

type PaymentController struct {
	pb.UnimplementedPaymentIntegrationServer
	service service.PaymentService
}

func NewPaymentController(service service.PaymentService) *PaymentController {
	return &PaymentController{service: service}
}

func (c *PaymentController) ProcessPayment(ctx context.Context, req *pb.Payment) (*pb.PaymentResponse, error) {
	p := &model.Payment{
		ID:     req.Id,
		Amount: req.Amount,
		From:   req.From,
		To:     req.To,
	}

	if err := c.service.ProcessPayment(p); err != nil {
		return &pb.PaymentResponse{Status: "error"}, err
	}

	return &pb.PaymentResponse{Status: "processed"}, nil
}
