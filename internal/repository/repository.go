package repository

import (
	"context"
	"pet-market/api"
	"pet-market/internal/models"
)

type UserRepository interface {
	Save(ctx context.Context, login string, password string) (int, error)
	GetUserByLogin(ctx context.Context, login string) (api.User, error)
}

type OrderRepository interface {
	Save(ctx context.Context, order models.Order, userID int) error
	GetAll(ctx context.Context, userID int) ([]models.Order, error)
	GetByOrderNumber(ctx context.Context, orderNum string) (models.Order, error)
}

type BalanceRepository interface {
	GetBalance(ctx context.Context, userID int) (api.Balance, error)
}

type WithdrawalsRepository interface {
	Save(ctx context.Context, orderNum string, sum float32, userID int) error
	GetAllByUserID(ctx context.Context, userID int) ([]models.Withdraw, error)
}
