package postgre

import (
	"context"
	"fmt"
	"service-courier/internal/config/dbcfg"
	"service-courier/observability/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func pingWithRetry(ctx context.Context, log logger.Logger, pool *pgxpool.Pool) error {
	var err error
	for i := 0; i < 3; i++ {
		err = pool.Ping(ctx)
		if err == nil {
			return err
		}
		log.Warn("database connection failed, retry...")
		time.Sleep(time.Second * 5)
	}
	return err
}

// InitPool - создание пула соединений с БД
func InitPool(ctx context.Context, log logger.Logger, env *dbcfg.DataBaseEnv) *pgxpool.Pool {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.DBUser,
		env.DBPass,
		env.DBHost,
		env.DBPort,
		env.DBName,
	)

	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Error("config database parsing error!")
	}

	cfg.MaxConns = 10
	cfg.MaxConnLifetime = time.Minute * 10
	cfg.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Error("pool creation error!")
	}
	if err := pingWithRetry(ctx, log, pool); err != nil {
		log.Error("database connection failed!")
	}
	return pool
}
