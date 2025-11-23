package courier

import (
	"errors"
	"fmt"
	"service-courier/internal/entity/courier"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func courierRepoMapError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return courier.ErrCourierExistPhone
		}
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return courier.ErrCourierNotFound
	}
	return fmt.Errorf("repo: failed to work with courier: %w", err)
}
