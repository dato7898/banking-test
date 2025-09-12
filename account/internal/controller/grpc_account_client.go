package controller

import (
	"account/internal/service"
	pb "account/proto"
	"context"
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
