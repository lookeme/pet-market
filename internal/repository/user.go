package repository

import (
	"context"
	"pet-market/api"
)

type UsrRepositoryImpl struct {
	pg *Postgres
}

func NewUsrRepository(pg *Postgres) *UsrRepositoryImpl {
	return &UsrRepositoryImpl{
		pg,
	}
}

func (r *UsrRepositoryImpl) Save(ctx context.Context, login string, password string) (int, error) {
	userID := 0
	err := r.pg.СonPool.QueryRow(
		ctx,
		"INSERT INTO users(login, pass) VALUES($1, $2) RETURNING id", login, password).Scan(&userID)
	if err != nil {
		return userID, err
	}
	return userID, nil
}
func (r *UsrRepositoryImpl) GetUserByLogin(ctx context.Context, login string) (api.User, error) {
	var usr api.User
	sqlStatement := `SELECT (id, login, pass) FROM users WHERE login = $1`
	err := r.pg.СonPool.QueryRow(ctx, sqlStatement, login).Scan(&usr)
	if err != nil {
		return api.User{}, err
	}
	return usr, nil
}
