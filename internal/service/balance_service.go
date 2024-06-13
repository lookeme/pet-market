package service

import (
	"context"
	"pet-market/api"
	"pet-market/internal/repository"
	"pet-market/internal/utils"
	"time"
)

type BalanceServiceIml struct {
	balanceRepo  repository.BalanceRepository
	withdrawRepo repository.WithdrawalsRepository
}

func NewBalanceService(balanceRepo repository.BalanceRepository,
	withdrawRepo repository.WithdrawalsRepository) *BalanceServiceIml {
	return &BalanceServiceIml{
		balanceRepo:  balanceRepo,
		withdrawRepo: withdrawRepo,
	}
}

func (b *BalanceServiceIml) GetBalance(ctx context.Context, userID int) (api.Balance, error) {
	return b.balanceRepo.GetBalance(ctx, userID)
}

func (b *BalanceServiceIml) AddWithdraw(ctx context.Context, userID int, withdraw api.RequestWithdraw) error {
	if !utils.VerifyLuhn(withdraw.Order) {
		return utils.ErrInvalidOrderNum
	}
	return b.withdrawRepo.Save(ctx, withdraw.Order, withdraw.Sum, userID)
}
func (b *BalanceServiceIml) GetAllWithdraws(ctx context.Context, userID int) ([]api.ResponseWithdraw, error) {
	var result []api.ResponseWithdraw
	withdrawals, err := b.withdrawRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return result, err
	}
	for _, w := range withdrawals {
		wr := api.ResponseWithdraw{
			Order:       w.OrderNum,
			Sum:         w.Sum,
			ProcessedAt: w.ProcessedAt.Format(time.RFC3339),
		}
		result = append(result, wr)
	}
	return result, nil
}
