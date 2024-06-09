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
	СonPool *pgxpool.Pool
	log     *logger.Logger
}

type Storage struct {
	UserRepository        UserRepository
	OrderRepository       OrderRepository
	BalanceRepository     BalanceRepository
	WithdrawalsRepository WithdrawalsRepository
}

func NewStorage(
	userRepo UserRepository,
	orderRepository OrderRepository,
	balanceRepository BalanceRepository,
	withdrawalsRepository WithdrawalsRepository,

) *Storage {
	return &Storage{
		UserRepository:        userRepo,
		OrderRepository:       orderRepository,
		BalanceRepository:     balanceRepository,
		WithdrawalsRepository: withdrawalsRepository,
	}
}

func (pg *Postgres) Close() error {
	pg.СonPool.Close()
	return nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.СonPool.Ping(ctx)
}

func New(ctx context.Context, log *logger.Logger, cfg *configuration.Storage) (*Postgres, error) {
	log.Log.Info("creating pool of conn to db...", zap.String("connString", cfg.ConnString))
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, cfg.ConnString)
		if err != nil {
			log.Log.Error(err.Error())
		}
		pgInstance = &Postgres{db, log}
	})
	err := db.StartMigration(pgInstance.СonPool)
	if err != nil {
		return nil, err
	}
	return pgInstance, nil
}
