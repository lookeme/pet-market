package service

import (
	"context"
	"pet-market/api"
)

type IUserService interface {
	CreateUser(ctx context.Context, user api.User) (int, error)
	GetUserByName(ctx context.Context, login string) (*api.User, error)
}

type IOrderService interface {
	CreateOrder(ctx context.Context, orderNum string, userID int) error
	GetUserOrders(ctx context.Context, userID int) ([]api.OrderResponse, error)
	GetOrder(ctx context.Context, orderNum string) (api.OrderResponse, error)
}

type IBalanceService interface {
	GetBalance(ctx context.Context, userID int) (api.Balance, error)
	AddWithdraw(ctx context.Context, userID int, withdraw api.RequestWithdraw) error
	GetAllWithdraws(ctx context.Context, userID int) ([]api.ResponseWithdraw, error)
}
