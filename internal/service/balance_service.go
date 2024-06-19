package service

import (
	"context"
	"pet-market/api"
	"pet-market/internal/repository"
	"pet-market/internal/utils"
	"time"
)

type BalanceService struct {
	balanceRepo  repository.IBalanceRepository
	withdrawRepo repository.IWithdrawalsRepository
}

func NewBalanceService(balanceRepo repository.IBalanceRepository,
	withdrawRepo repository.IWithdrawalsRepository) *BalanceService {
	return &BalanceService{
		balanceRepo:  balanceRepo,
		withdrawRepo: withdrawRepo,
	}
}

func (b *BalanceService) GetBalance(ctx context.Context, userID int) (api.Balance, error) {
	balance, err := b.balanceRepo.GetBalance(ctx, userID)
	if err != nil {
		return api.Balance{}, err
	}
	balance.Current = balance.Current - balance.Withdrawn
	return balance, nil
}

func (b *BalanceService) AddWithdraw(ctx context.Context, userID int, withdraw api.RequestWithdraw) error {
	if !utils.VerifyLuhn(withdraw.Order) {
		return utils.ErrInvalidOrderNum
	}
	return b.withdrawRepo.Save(ctx, withdraw.Order, withdraw.Sum, userID)
}
func (b *BalanceService) GetAllWithdraws(ctx context.Context, userID int) ([]api.ResponseWithdraw, error) {
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
