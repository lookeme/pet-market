package repository

import (
	"context"
	"pet-market/api"
	"pet-market/internal/models"
	"pet-market/internal/utils"

	"github.com/jackc/pgx/v5"
)

type WithdrawRepositoryImpl struct {
	pg *Postgres
}

func NewWithdrawRepository(pg *Postgres) *WithdrawRepositoryImpl {
	return &WithdrawRepositoryImpl{
		pg,
	}
}

func (w *WithdrawRepositoryImpl) GetAllByUserID(ctx context.Context, userID int) ([]models.Withdraw, error) {
	rows, err := w.pg.СonPool.Query(ctx, "SELECT order_num, processed_at, sum, user_id FROM withdrawals WHERE user_id = $1", userID)
	if err != nil {
		return []models.Withdraw{}, nil
	}
	defer rows.Close()
	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Withdraw])
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (w *WithdrawRepositoryImpl) Save(ctx context.Context, orderNum string, sum float32, userID int) error {
	tx, err := w.pg.СonPool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	if err != nil {
		return err
	}
	var balance api.Balance
	sqlStatement := "SELECT Coalesce(current, 0) as current,  Coalesce(withdrawn, 0) as withdrawn FROM (SELECT SUM(accrual) as current, user_id FROM orders WHERE user_id = $1 group by user_id) T1 LEFT JOIN (SELECT SUM(sum) as withdrawn, user_id FROM withdrawals WHERE user_id = $1 group by user_id) T2 ON T1.user_id = T2.user_id;"
	err = tx.QueryRow(ctx, sqlStatement, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return err
	}
	if balance.Current-(sum+balance.Withdrawn) < 0 {
		return utils.ErrInsufficientFunds
	}
	_, err = tx.Exec(ctx, "INSERT INTO withdrawals (order_num, sum, user_id) VALUES ($1, $2, $3 )", orderNum, sum, userID)
	if err != nil {
		return err
	}
	return nil
}
