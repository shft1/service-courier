package courierdb

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"service-courier/internal/domain/courier"
)

const tableCouriers = "couriers"

// courierRepository - репозиторий курьера
type courierRepository struct {
	pool         *pgxpool.Pool
	queryBuilder sq.StatementBuilderType
	txManager    txManagerGetConnection
}

// NewCourierRepository - конструктор репозитория курьера
func NewCourierRepository(pool *pgxpool.Pool, txManager txManagerGetConnection) *courierRepository {
	return &courierRepository{
		pool:         pool,
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		txManager:    txManager,
	}
}

// Create - создать нового курьера в БД
func (cr *courierRepository) Create(ctx context.Context, cour *courier.CourierCreate) (int64, error) {
	courRow := domainToRowCreate(cour)
	var id int64

	columns, values := []string{"name", "phone"}, []any{courRow.Name, courRow.Phone}
	if courRow.TransportType != nil {
		columns = append(columns, "transport_type")
		values = append(values, courRow.TransportType)
	}
	if courRow.Status != nil {
		columns = append(columns, "status")
		values = append(values, courRow.Status)
	}
	queryCreate := cr.queryBuilder.
		Insert(tableCouriers).
		Columns(columns...).
		Values(values...).
		Suffix("RETURNING id")

	query, args, err := queryCreate.ToSql()
	if err != nil {
		return -1, fmt.Errorf("repo: failed to build query: %w", err)
	}
	if err = cr.pool.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		return -1, mapError(err)
	}
	return id, nil
}

// Update - обновить курьера в БД
func (cr *courierRepository) Update(ctx context.Context, cour *courier.CourierUpdate) (int64, error) {
	courRow := domainToRowUpdate(cour)
	var id int64

	query := `
	UPDATE couriers
	SET name = COALESCE($1, name),
		phone = COALESCE($2, phone),
		status = COALESCE($3, status),
		transport_type = COALESCE($4, transport_type),
		updated_at = now()
	WHERE id = $5
	RETURNING id;`
	args := []any{courRow.Name, courRow.Phone, courRow.Status, courRow.TransportType, courRow.ID}

	if err := cr.pool.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		return -1, mapError(err)
	}
	return id, nil
}

// GetByID - получить курьера по ID из БД
func (cr *courierRepository) GetByID(ctx context.Context, id int64) (*courier.Courier, error) {
	var courRow courierRow

	query := `
	SELECT id, name, phone, status, transport_type
	FROM couriers
	WHERE id = $1;`
	args := []any{id}

	err := cr.pool.QueryRow(
		ctx, query, args...,
	).Scan(&courRow.ID, &courRow.Name, &courRow.Phone, &courRow.Status, &courRow.TransportType)

	if err != nil {
		return nil, mapError(err)
	}
	return rowToDomainCourier(&courRow), nil
}

// GetMulti - получить курьеров из БД
func (cr *courierRepository) GetMulti(ctx context.Context) ([]courier.Courier, error) {
	var coursRow []courier.Courier

	query := `
	SELECT id, name, phone, status, transport_type, created_at, updated_at
	FROM couriers;`

	rows, err := cr.pool.Query(ctx, query)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()
	for rows.Next() {
		var courRow courierRow
		if err := rows.Scan(
			&courRow.ID,
			&courRow.Name,
			&courRow.Phone,
			&courRow.Status,
			&courRow.TransportType,
			&courRow.CreatedAt,
			&courRow.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("repo: failed to scan data: %w", err)
		}
		coursRow = append(coursRow, *rowToDomainCourier(&courRow))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo: failed while reading: %w", err)
	}
	return coursRow, nil
}

// GetAvailable - получить курьера с наименьшим кол-вом доставок
func (cr *courierRepository) GetAvailable(ctx context.Context) (*courier.Courier, error) {
	var courRow courierRow

	tx, err := cr.txManager.GetTx(ctx)
	if err != nil {
		return nil, err
	}
	query := `
	WITH candidate AS (
		SELECT c.id as cand_id
		FROM couriers AS c
		LEFT JOIN delivery as d on c.id = d.courier_id
		WHERE c.status = 'available'
		GROUP BY c.id
		ORDER BY COUNT(d.id)
		LIMIT 1
	)
	SELECT c.id, c.name, c.phone, c.status, c.transport_type, c.created_at, c.updated_at
	FROM couriers AS c
	JOIN candidate ON c.id = candidate.cand_id
	FOR UPDATE SKIP LOCKED;`

	err = tx.QueryRow(
		ctx, query,
	).Scan(
		&courRow.ID,
		&courRow.Name,
		&courRow.Phone,
		&courRow.Status,
		&courRow.TransportType,
		&courRow.CreatedAt,
		&courRow.UpdatedAt,
	)
	if err != nil {
		return nil, courier.ErrCourierAvailable
	}
	return rowToDomainCourier(&courRow), nil
}

// SetBusy - установить статус курьера на занятого по ID
func (cr *courierRepository) SetBusy(ctx context.Context, id int64) (int64, error) {
	var courID int64

	tx, err := cr.txManager.GetTx(ctx)
	if err != nil {
		return -1, mapError(err)
	}
	query := `
	UPDATE couriers
	SET status = 'busy'
	WHERE id = $1
	RETURNING id;`
	args := []any{id}

	if err := tx.QueryRow(ctx, query, args...).Scan(&courID); err != nil {
		return -1, mapError(err)
	}
	return courID, nil
}

// SetAvailable - установить статус курьера на свободного по ID
func (cr *courierRepository) SetAvailable(ctx context.Context, id int64) (int64, error) {
	var courID int64

	tx, err := cr.txManager.GetTx(ctx)
	if err != nil {
		return -1, mapError(err)
	}
	query := `
	UPDATE couriers
	SET status = 'available'
	WHERE id = $1
	RETURNING id;`
	args := []any{id}

	if err := tx.QueryRow(ctx, query, args...).Scan(&courID); err != nil {
		return -1, mapError(err)
	}
	return courID, nil
}

// ReleaseStaleBusy - освободить курьеров, выполнивших заказ
func (cr *courierRepository) ReleaseStaleBusy(ctx context.Context) error {
	query := `
	UPDATE couriers
	SET status = 'available'
	WHERE status = 'busy' AND id NOT IN (
		SELECT courier_id
		FROM delivery
		WHERE NOW() <= deadline
	);`
	if _, err := cr.pool.Exec(ctx, query); err != nil {
		return mapError(err)
	}
	return nil
}
