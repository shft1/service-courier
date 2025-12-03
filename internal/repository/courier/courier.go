package courier

import (
	"context"
	"fmt"
	"service-courier/internal/entity/courier"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const tableCouriers = "couriers"

type courierRepository struct {
	pool        *pgxpool.Pool
	queryBilder sq.StatementBuilderType
	txManager   txManagerGetConnection
}

func NewCourierRepository(pool *pgxpool.Pool, txManager txManagerGetConnection) *courierRepository {
	return &courierRepository{
		pool:        pool,
		queryBilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		txManager:   txManager,
	}
}

func (cr *courierRepository) Create(ctx context.Context, c *courier.CourierCreate) error {
	columns, values := []string{"name", "phone"}, []any{c.Name, c.Phone}
	if c.TransportType != nil {
		columns = append(columns, "transport_type")
		values = append(values, c.TransportType)
	}
	if c.Status != nil {
		columns = append(columns, "status")
		values = append(values, c.Status)
	}
	queryCreate := cr.queryBilder.
		Insert(tableCouriers).
		Columns(columns...).
		Values(values...)

	query, args, err := queryCreate.ToSql()
	if err != nil {
		return fmt.Errorf("repo: failed to build query: %w", err)
	}

	_, err = cr.pool.Exec(ctx, query, args...)

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
		transport_type = COALESCE($4, transport_type),
		updated_at = now()
	WHERE id = $5;`
	args := []any{c.Name, c.Phone, c.Status, c.TransportType, c.ID}

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
	SELECT id, name, phone, status, transport_type
	FROM couriers
	WHERE id = $1;`
	args := []any{id}

	err := cr.pool.QueryRow(
		ctx, query, args...,
	).Scan(&c.ID, &c.Name, &c.Phone, &c.Status, &c.TransportType)

	if err != nil {
		return nil, courierRepoMapError(err)
	}
	return &c, nil
}

func (cr *courierRepository) GetMulti(ctx context.Context) ([]courier.CourierGet, error) {
	var couriers []courier.CourierGet
	query := `SELECT id, name, phone, status, transport_type FROM couriers;`

	rows, err := cr.pool.Query(ctx, query)

	if err != nil {
		return nil, courierRepoMapError(err)
	}
	defer rows.Close()
	for rows.Next() {
		var c courier.CourierGet
		if err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Status, &c.TransportType); err != nil {
			return nil, fmt.Errorf("repo: failed to scan data: %w", err)
		}
		couriers = append(couriers, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo: failed while reading: %w", err)
	}
	return couriers, nil
}

func (cr *courierRepository) GetAvailable(ctx context.Context) (*courier.CourierGet, error) {
	var c courier.CourierGet

	tx, err := cr.txManager.GetTx(ctx)
	if err != nil {
		return nil, err
	}
	query := `
	WITH candidate AS (
		SELECT c.id as cand_id
		FROM couriers as c
		LEFT JOIN delivery as d on c.id = d.courier_id
		WHERE c.status = 'available'
		GROUP BY c.id
		ORDER BY COUNT(d.id)
		LIMIT 1
	)
	SELECT couriers.id, couriers.name, couriers.phone, couriers.status, couriers.transport_type
	FROM couriers
	JOIN candidate ON couriers.id = candidate.cand_id
	FOR UPDATE SKIP LOCKED;`

	err = tx.QueryRow(
		ctx, query,
	).Scan(&c.ID, &c.Name, &c.Phone, &c.Status, &c.TransportType)
	if err != nil {
		return nil, courier.ErrCourierAvailable
	}
	return &c, nil
}

func (cr *courierRepository) SetBusy(ctx context.Context, courierID int) error {
	tx, err := cr.txManager.GetTx(ctx)
	if err != nil {
		return courierRepoMapError(err)
	}

	query := `
	UPDATE couriers
	SET status = 'busy'
	WHERE id = $1;`
	args := []any{courierID}

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return courierRepoMapError(err)
	}
	if result.RowsAffected() == 0 {
		return courier.ErrCourierNotFound
	}
	return nil
}

func (cr *courierRepository) SetAvailable(ctx context.Context, courierID int) error {
	tx, err := cr.txManager.GetTx(ctx)
	if err != nil {
		return courierRepoMapError(err)
	}

	query := `
	UPDATE couriers
	SET status = 'available'
	WHERE id = $1;`
	args := []any{courierID}

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return courierRepoMapError(err)
	}
	if result.RowsAffected() == 0 {
		return courier.ErrCourierNotFound
	}
	return nil
}
