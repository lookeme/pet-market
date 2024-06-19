package repository

import (
	"context"
	"pet-market/internal/configuration"
	"pet-market/internal/logger"
	"pet-market/internal/repository/db"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

type Postgres struct {
	ConPool *pgxpool.Pool
	log     *logger.Logger
}

type Storage struct {
	UserRepository        IUserRepository
	OrderRepository       IOrderRepository
	BalanceRepository     IBalanceRepository
	WithdrawalsRepository IWithdrawalsRepository
}

func (pg *Postgres) Close() error {
	pg.ConPool.Close()
	return nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.ConPool.Ping(ctx)
}

func New(ctx context.Context, log *logger.Logger, cfg *configuration.Storage) (*Postgres, error) {
	log.Log.Info("creating pool of conn to db...", zap.String("connString", cfg.ConnString))
	pgOnce.Do(func() {
		conPool, err := pgxpool.New(ctx, cfg.ConnString)
		if err != nil {
			log.Log.Error(err.Error())
		}
		pgInstance = &Postgres{conPool, log}
	})
	err := db.StartMigration(pgInstance.ConPool)
	if err != nil {
		return nil, err
	}
	return pgInstance, nil
}
