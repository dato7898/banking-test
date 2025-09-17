package controller

import (
	"account/internal/service"
	pb "account/proto"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

type AccountController struct {
	pb.UnimplementedAccountClientServer
	accountService service.AccountService
}

func NewAccountController(accountService service.AccountService) *AccountController {
	return &AccountController{accountService: accountService}
}

func (c *AccountController) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.Account, error) {
	iban := req.Iban

	a, err := c.accountService.GetAccount(iban)
	if err != nil {
		return nil, err
	}

	return &pb.Account{
		Id:     a.ID,
		Iban:   a.Iban,
		Amount: a.Amount,
		UserID: int32(a.UserID),
	}, nil
}

func (c *AccountController) Replenishment(ctx context.Context, req *pb.OperationRequest) (*emptypb.Empty, error) {
	if err := c.accountService.Replenishment(req.Iban, req.Amount); err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (c *AccountController) Withdrawal(ctx context.Context, req *pb.OperationRequest) (*emptypb.Empty, error) {
	if err := c.accountService.Withdrawal(req.Iban, req.Amount); err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
