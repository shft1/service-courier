package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TxManagerGetConnection interface {
	GetTx(ctx context.Context) (pgx.Tx, error)
}
