package service

import (
	"context"
	"pet-market/api"
)

type UserService interface {
	CreateUser(ctx context.Context, user api.User) error
	GetUserByName(ctx context.Context, login string) (*api.User, error)
}

type OrderService interface {
	CreateOrder(ctx context.Context, orderNum string, userID int) error
	GetUserOrders(ctx context.Context, userID int) ([]api.OrderResponse, error)
	GetOrder(ctx context.Context, orderNum string) (api.OrderResponse, error)
}

type BalanceService interface {
	GetBalance(ctx context.Context, userID int) (api.Balance, error)
	AddWithdraw(ctx context.Context, userID int, withdraw api.RequestWithdraw) error
	GetAllWithdraws(ctx context.Context, userID int) ([]api.ResponseWithdraw, error)
}
