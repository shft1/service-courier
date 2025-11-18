package postgre

import (
	"context"
	"fmt"
	"service-courier/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func pingWithRetry(ctx context.Context, pool *pgxpool.Pool) error {
	var err error
	for i := 0; i < 3; i++ {
		err = pool.Ping(ctx)
		if err == nil {
			return err
		}
		fmt.Println("database connection failed, retry...")
		time.Sleep(time.Second * 5)
	}
	return err
}

func InitPool(ctx context.Context, env *config.Env) *pgxpool.Pool {
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
		fmt.Println("config database parsing error!")
	}

	cfg.MaxConns = 10
	cfg.MaxConnLifetime = time.Minute * 10
	cfg.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		fmt.Println("pool creation error!")
	}
	if err := pingWithRetry(ctx, pool); err != nil {
		fmt.Println("database connection failed!")
	}
	return pool
}
