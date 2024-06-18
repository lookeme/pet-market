package service

import (
	"context"
	"errors"
	"pet-market/api"
	"pet-market/internal/integration"
	"pet-market/internal/models"
	"pet-market/internal/repository"
	"pet-market/internal/utils"
	"time"

	"github.com/jackc/pgx/v5"
)

type OrderServiceIml struct {
	accural   integration.AccuralIntegration
	orderRepo repository.OrderRepository
}

func NewOrderService(accural *integration.AccuralIntegration,
	orderRepo repository.OrderRepository) *OrderServiceIml {
	return &OrderServiceIml{
		accural:   *accural,
		orderRepo: orderRepo,
	}
}
func (i *OrderServiceIml) CreateOrder(ctx context.Context, orderNum string, userID int) error {
	if !utils.VerifyLuhn(orderNum) {
		return utils.ErrInvalidOrderNum
	}
	existedOrder, errPg := i.orderRepo.GetByOrderNumber(ctx, orderNum)
	if errors.Is(errPg, pgx.ErrNoRows) {
		order, getError := i.accural.GetOrder(orderNum)
		if getError != nil {
			return getError
		}
		orderToSave := models.Order{
			Accrual: order.Accrual,
			OrderID: order.Order,
			Status:  order.Status,
			UserID:  userID,
		}
		return i.orderRepo.Save(ctx, orderToSave, userID)
	} else if existedOrder.UserID == userID {
		return utils.ErrOrderNumAttachedToAnotherUser
	} else {
		return utils.ErrOrderNumIsAlreadyRegistered
	}
}

func (i *OrderServiceIml) GetUserOrders(ctx context.Context, userID int) ([]api.OrderResponse, error) {
	var result []api.OrderResponse
	orders, err := i.orderRepo.GetAll(ctx, userID)
	if err != nil {
		return result, err
	}
	for _, order := range orders {
		o := api.OrderResponse{
			Number:     order.OrderID,
			Accrual:    order.Accrual,
			Status:     order.Status,
			UploadedAt: order.UploadedAt.Format(time.RFC3339),
		}
		result = append(result, o)
	}
	return result, nil
}

func (i *OrderServiceIml) GetOrder(ctx context.Context, orderNum string) (api.OrderResponse, error) {
	order, err := i.orderRepo.GetByOrderNumber(ctx, orderNum)
	if err != nil {
		return api.OrderResponse{}, err
	}
	return api.OrderResponse{
		Number:     order.OrderID,
		Accrual:    order.Accrual,
		Status:     order.Status,
		UploadedAt: order.UploadedAt.Format(time.RFC3339),
	}, nil
}
