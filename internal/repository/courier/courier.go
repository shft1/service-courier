package courier

import (
	"context"
	"errors"
	"fmt"
	"service-courier/internal/entity/courier"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func courierRepoMapError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return courier.ErrCourierExistPhone
		default:
			return courier.ErrDatabase
		}
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return courier.ErrCourierNotFound
	}
	return fmt.Errorf("repo: failed to work with courier: %w", err)
}

type courierRepository struct {
	pool *pgxpool.Pool
}

func NewCourierRepository(pool *pgxpool.Pool) *courierRepository {
	return &courierRepository{
		pool: pool,
	}
}

func (cr *courierRepository) Create(ctx context.Context, c *courier.CourierCreate) error {
	query := `INSERT INTO couriers (name, phone, status) VALUES ($1, $2, $3);`
	args := []any{c.Name, c.Phone, c.Status}
	_, err := cr.pool.Exec(ctx, query, args...)
	if err != nil {
		return courierRepoMapError(err)
	}
	return nil
}

func (cr *courierRepository) Update(ctx context.Context, c *courier.CourierUpdate) error {
	query := `
	UPDATE couriers
	SET name = COALESCE($1, name),
		phone = COALESCE($2, phone),
		status = COALESCE($3, status),
		updated_at = now()
	WHERE id = $4;`
	args := []any{c.Name, c.Phone, c.Status, c.ID}
	result, err := cr.pool.Exec(ctx, query, args...)
	if err != nil {
		return courierRepoMapError(err)
	}
	if result.RowsAffected() == 0 {
		return courier.ErrCourierNotFound
	}
	return nil
}

func (cr *courierRepository) GetByID(ctx context.Context, id int) (*courier.CourierGet, error) {
	var c courier.CourierGet
	query := `
	SELECT id, name, phone, status
	FROM couriers
	WHERE id = $1;`
	args := []any{id}
	err := cr.pool.QueryRow(
		ctx, query, args...,
	).Scan(&c.ID, &c.Name, &c.Phone, &c.Status)
	if err != nil {
		return nil, courierRepoMapError(err)
	}
	return &c, nil
}

func (cr *courierRepository) GetMulti(ctx context.Context) ([]courier.CourierGet, error) {
	var couriers []courier.CourierGet
	query := `SELECT id, name, phone, status FROM couriers;`
	rows, err := cr.pool.Query(ctx, query)
	if err != nil {
		return nil, courierRepoMapError(err)
	}
	defer rows.Close()
	for rows.Next() {
		var c courier.CourierGet
		if err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Status); err != nil {
			return nil, fmt.Errorf("repo: failed to scan data: %w", err)
		}
		couriers = append(couriers, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo: failed while reading: %w", err)
	}
	return couriers, nil
}
