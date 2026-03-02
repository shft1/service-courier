package deliverydb

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/shft1/service-courier/internal/domain/delivery"
)

// mapError - маппинг ошибок репозитория доставок
func mapError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return delivery.ErrDeliveryExist
		}
		if pgErr.Code == "23503" {
			return delivery.ErrDeliveryInvalidAssignCourier
		}
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return delivery.ErrDeliveryNotFound
	}
	return fmt.Errorf("repo: failed to work with delivery: %w", err)
}
