package postgre

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txContextKey struct{}

// txManagerPostgre - транзакционный менеджер
type txManagerPostgre struct {
	pool *pgxpool.Pool
}

// NewTxManagerPostgre - конструктор транзакционного менеджера
func NewTxManagerPostgre(pool *pgxpool.Pool) *txManagerPostgre {
	return &txManagerPostgre{pool: pool}
}

// Do - обернуть переданную функцию в транзакцию
func (tm *txManagerPostgre) Do(prnt context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.begin(prnt)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tm.rollback(prnt, tx)

	ctxTx := context.WithValue(prnt, txContextKey{}, tx)
	if err := fn(ctxTx); err != nil {
		tm.rollback(prnt, tx)
		return fmt.Errorf("failed to execute transaction: %w", err)
	}

	if err := tm.commit(prnt, tx); errors.Is(err, pgx.ErrTxCommitRollback) {
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
