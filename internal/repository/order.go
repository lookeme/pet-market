package repository

import (
	"context"
	"pet-market/internal/models"

	"github.com/jackc/pgx/v5"
)

type OrderRepositoryImpl struct {
	pg *Postgres
}

func NewOrderRepository(pg *Postgres) *OrderRepositoryImpl {
	return &OrderRepositoryImpl{
		pg,
	}
}

func (r *OrderRepositoryImpl) Save(ctx context.Context, order models.Order, userID int) error {
	_, err := r.pg.СonPool.Exec(
		ctx, "INSERT INTO orders(order_id,status, accrual, user_id) VALUES($1, $2, $3, $4)",
		order.OrderID, order.Status, order.Accrual, userID)
	return err
}

func (r *OrderRepositoryImpl) GetAll(ctx context.Context, userID int) ([]models.Order, error) {
	rows, err := r.pg.СonPool.Query(ctx, "SELECT order_id, accrual, status, uploaded_at, user_id FROM orders WHERE user_id = $1", userID)
	if err != nil {
		return []models.Order{}, nil
	}
	defer rows.Close()
	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Order])
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *OrderRepositoryImpl) GetByOrderNumber(ctx context.Context, orderNum string) (models.Order, error) {
	var order models.Order
	sqlStatement := `SELECT order_id, accrual, status, uploaded_at,  user_id  FROM orders WHERE order_id = $1`
	err := r.pg.СonPool.QueryRow(ctx, sqlStatement, orderNum).Scan(&order.OrderID, &order.Accrual, &order.Status, &order.UploadedAt, &order.UserID)
	return order, err
}
