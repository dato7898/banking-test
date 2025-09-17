package repository

import (
	"context"
	pb "payment/proto"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AccountClient interface {
	GetAccount(req *pb.GetAccountRequest) (*pb.Account, error)
	Replenishment(req *pb.OperationRequest) (*emptypb.Empty, error)
	Withdrawal(req *pb.OperationRequest) (*emptypb.Empty, error)
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

func (c *accountClient) Replenishment(req *pb.OperationRequest) (*emptypb.Empty, error) {
	return c.client.Replenishment(context.Background(), req)
}

func (c *accountClient) Withdrawal(req *pb.OperationRequest) (*emptypb.Empty, error) {
	return c.client.Withdrawal(context.Background(), req)
}
