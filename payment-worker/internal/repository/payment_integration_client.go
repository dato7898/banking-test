package repository

import (
	"context"
	pb "payment-worker/proto"

	"google.golang.org/grpc"
)

type PaymentIntegrationClient interface {
	ProcessPayment(p *pb.Payment) error
}

type paymentIntegrationClient struct {
	client pb.PaymentIntegrationClient
}

func NewPaymentIntegrationClient(conn *grpc.ClientConn) PaymentIntegrationClient {
	return &paymentIntegrationClient{client: pb.NewPaymentIntegrationClient(conn)}
}

func (c *paymentIntegrationClient) ProcessPayment(p *pb.Payment) error {
	_, err := c.client.ProcessPayment(context.Background(), p)
	return err
}
