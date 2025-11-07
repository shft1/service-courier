package repository

import (
	"context"
	"errors"
	"fmt"
	"service-courier/internal/entity/courier"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type courierRepository struct {
	pool *pgxpool.Pool
}

func NewCourierRepository(pool *pgxpool.Pool) courier.CourierRepository {
	return &courierRepository{
		pool: pool,
	}
}

func (cr *courierRepository) Create(c *courier.CourierCreate) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	queryBuilder := psql.
		Insert("couriers").
		Columns("name", "phone", "status").
		Values(c.Name, c.Phone, c.Status)
	query, args, _ := queryBuilder.ToSql()
	_, err := cr.pool.Exec(context.Background(), query, args...)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "couriers_phone_key"):
			return courier.ErrCourierExistPhone
		default:
			return fmt.Errorf("repo: failed to create courier: %w", err)
		}
	}
	return nil
}

func (cr *courierRepository) Update(c *courier.CourierUpdate) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	queryBuilder := psql.
		Update("couriers").
		Set("name", sq.Expr("COALESCE(?, name)", c.Name)).
		Set("phone", sq.Expr("COALESCE(?, phone)", c.Phone)).
		Set("status", sq.Expr("COALESCE(?, status)", c.Status)).
		Where(sq.Eq{"id": c.ID})
	query, args, _ := queryBuilder.ToSql()
	result, err := cr.pool.Exec(context.Background(), query, args...)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "couriers_phone_key"):
			return courier.ErrCourierExistPhone
		default:
			return fmt.Errorf("repo: failed to update courier: %w", err)
		}
	}
	if result.RowsAffected() == 0 {
		return courier.ErrCourierNotFound
	}
	return nil
}

func (cr *courierRepository) GetByID(id int) (*courier.CourierGet, error) {
	var c courier.CourierGet
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	queryBuilder := psql.
		Select("id", "name", "phone", "status").
		From("couriers").
		Where(sq.Eq{"id": id})
	query, args, _ := queryBuilder.ToSql()
	err := cr.pool.QueryRow(
		context.Background(), query, args...,
	).Scan(&c.ID, &c.Name, &c.Phone, &c.Status)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, courier.ErrCourierNotFound
		default:
			return nil, fmt.Errorf("repo: failed to get courier by id: %w", err)
		}
	}
	return &c, nil
}

func (cr *courierRepository) GetMulti() ([]courier.CourierGet, error) {
	var couriers []courier.CourierGet
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	queryBuilder := psql.
		Select("id", "name", "phone", "status").
		From("couriers")
	query, _, _ := queryBuilder.ToSql()
	rows, err := cr.pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("repo: failed to get couriers: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var c courier.CourierGet
		if err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Status); err != nil {
			return nil, fmt.Errorf("repo: failed to format data: %w", err)
		}
		couriers = append(couriers, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo: failed while reading: %w", err)
	}
	return couriers, nil
}
