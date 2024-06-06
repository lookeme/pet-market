package repository

import (
	"context"
	"pet-market/api"
	"pet-market/internal/repository/db"
)

type UsrRepositoryImpl struct {
	pg db.Postgres
}

func (r *UsrRepositoryImpl) Save(userName string, password string) (int, error) {
	lastInsertID := 0
	err := r.pg.СonPool.QueryRow(
		context.Background(),
		"INSERT INTO users(login, pass) VALUES($1, $2) RETURNING id", userName, password).Scan(&lastInsertID)
	if err != nil {
		return lastInsertID, err
	}
	return lastInsertID, nil
}
func (r *UsrRepositoryImpl) GetUserByName(userName string) (api.User, error) {
	var usr api.User
	sqlStatement := `SELECT (login, pass) FROM users WHERE login = $1`
	err := r.pg.СonPool.QueryRow(context.Background(), sqlStatement, userName).Scan(&usr)
	if err != nil {
		return api.User{}, err
	}
	return usr, nil
}
