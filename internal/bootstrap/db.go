package bootstrap

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPool(ctx context.Context, env *Env) *pgxpool.Pool {
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
		fmt.Println("Ошибка парсинга конфига !")
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		fmt.Println("Ошибка создания пула !")
	}
	err = pool.Ping(ctx)
	if err != nil {
		fmt.Println("Ошибка работы пула !")
	}
	return pool
}
