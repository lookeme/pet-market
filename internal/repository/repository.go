package repository

import "pet-market/api"

type UserRepository interface {
	Save(userName string, password string) (int, error)
	GetUserByName(userName string) (api.User, error)
}

type OrderRepository interface {
	Save(orderNum string, userId int) error
	GetAll(userId string) ([]api.Order, error)
}

type BalanceRepository interface {
	Save(orderNum string, userId int)
	GetBalance(userId int) (api.Balance, error)
}

type WithdrawalsRepository interface {
	Save(orderNum string, userId int, sum int) error
	GetAllByUserId(userId int) ([]api.Withdraw, error)
}
