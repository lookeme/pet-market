package repository

import "pet-market/api"

type OrderRepositoryImpl struct {
	pg *Postgres
}

func (r *OrderRepositoryImpl) Save(orderNum string, userId int) error {
	return nil
}

func (r *OrderRepositoryImpl) GetAll(userId string) ([]api.Order, error) {
	return nil, nil
}
