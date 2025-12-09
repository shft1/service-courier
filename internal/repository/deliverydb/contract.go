package deliverydb

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type txManagerGetConnection interface {
	GetTx(ctx context.Context) (pgx.Tx, error)
}
