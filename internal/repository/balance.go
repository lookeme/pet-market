package repository

import (
	"context"
	"pet-market/api"
)

type BalanceRepositoryImpl struct {
	pg *Postgres
}

func NewBalanceRepository(pg *Postgres) *BalanceRepositoryImpl {
	return &BalanceRepositoryImpl{
		pg,
	}
}

func (b *BalanceRepositoryImpl) GetBalance(ctx context.Context, userID int) (api.Balance, error) {
	var balance api.Balance
	sqlStatement := "SELECT Coalesce(current, 0) as current,  Coalesce(withdrawn, 0) as withdrawn FROM (SELECT SUM(accrual) as current, user_id FROM orders WHERE user_id = $1 group by user_id) T1 LEFT JOIN (SELECT SUM(sum) as withdrawn, user_id FROM withdrawals WHERE user_id = $1 group by user_id) T2 ON T1.user_id = T2.user_id;"
	err := b.pg.СonPool.QueryRow(ctx, sqlStatement, userID).Scan(&balance.Current, &balance.Withdrawn)
	return balance, err
}
