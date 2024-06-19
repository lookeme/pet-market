package repository

import (
	"context"
	"pet-market/api"
	"pet-market/internal/models"
)

type IUserRepository interface {
	Save(ctx context.Context, login string, password string) (int, error)
	GetUserByLogin(ctx context.Context, login string) (api.User, error)
}

type IOrderRepository interface {
	Save(ctx context.Context, order models.Order, userID int) error
	GetAll(ctx context.Context, userID int) ([]models.Order, error)
	GetByOrderNumber(ctx context.Context, orderNum string) (models.Order, error)
}

type IBalanceRepository interface {
	GetBalance(ctx context.Context, userID int) (api.Balance, error)
}

type IWithdrawalsRepository interface {
	Save(ctx context.Context, orderNum string, sum float32, userID int) error
	GetAllByUserID(ctx context.Context, userID int) ([]models.Withdraw, error)
}
