package repository

import (
	"context"
	pb "payment/proto"

	"google.golang.org/grpc"
)

type AccountClient interface {
	GetAccount(iban *pb.GetAccountRequest) (*pb.Account, error)
}

type accountClient struct {
	client pb.AccountClientClient
}

func NewAccountClient(conn *grpc.ClientConn) AccountClient {
	return &accountClient{client: pb.NewAccountClientClient(conn)}
}

func (c *accountClient) GetAccount(req *pb.GetAccountRequest) (*pb.Account, error) {
	return c.client.GetAccount(context.Background(), req)
}
