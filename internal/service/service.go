package service

import (
	"pet-market/api"
)

type UserService interface {
	CreateUser(user api.User) error
	GetUserByName(userName string) (*api.User, error)
}

type OrderService interface {
	CreateOrder(orderNum string, userName string) error
	GetUserOrders(userName string) ([]api.Order, error)
}

type BalanceService interface {
	GetBalance(userName string) (api.Balance, error)
	AddWithdraw(userName string, sum int) error
	GetAllWithdraw(userName string) ([]api.Withdraw, error)
}
