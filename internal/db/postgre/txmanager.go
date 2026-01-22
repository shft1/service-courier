package postgre

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"service-courier/observability/logger"
)

type txContextKey struct{}

// txManagerPostgre - транзакционный менеджер
type txManagerPostgre struct {
	log  logger.Logger
	pool *pgxpool.Pool
}

// NewTxManagerPostgre - конструктор транзакционного менеджера
func NewTxManagerPostgre(log logger.Logger, pool *pgxpool.Pool) *txManagerPostgre {
	return &txManagerPostgre{
		log:  log,
		pool: pool,
	}
}

// Do - обернуть переданную функцию в транзакцию
func (tm *txManagerPostgre) Do(prnt context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.begin(prnt)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err := tm.rollback(prnt, tx); err != nil {
			tm.log.Warn("failed to rollback transaction gracefully", logger.NewField("error", err))
		}
	}()

	ctxTx := context.WithValue(prnt, txContextKey{}, tx)

	err = fn(ctxTx)
	if err != nil {
		txErr := tm.rollback(prnt, tx)
		if txErr != nil {
			tm.log.Warn("failed to rollback transaction gracefully", logger.NewField("error", err))
		}
		return fmt.Errorf("failed to execute transaction: %w", err)
	}

	err = tm.commit(prnt, tx)
	if errors.Is(err, pgx.ErrTxCommitRollback) {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (tm *txManagerPostgre) begin(prnt context.Context) (pgx.Tx, error) {
	return tm.pool.Begin(prnt)
}

func (tm *txManagerPostgre) rollback(prnt context.Context, tx pgx.Tx) error {
	return tx.Rollback(prnt)
}

func (tm *txManagerPostgre) commit(prnt context.Context, tx pgx.Tx) error {
	return tx.Commit(prnt)
}

// GetTx - получить транзакционное соединение из контекста
func (tm *txManagerPostgre) GetTx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(txContextKey{}).(pgx.Tx)
	if !ok {
		return nil, errors.New("failed to get transaction")
	}
	return tx, nil
}
