package deliverydb

import (
	"context"
	"service-courier/internal/domain/delivery"
	"service-courier/internal/domain/order"

	"github.com/jackc/pgx/v5/pgxpool"
)

// deliveryRepository - репозиторий доставки
type deliveryRepository struct {
	pool      *pgxpool.Pool
	txManager txManagerGetConnection
}

// NewDeliveryRepository - конструктор репозитория доставки
func NewDeliveryRepository(pool *pgxpool.Pool, txManager txManagerGetConnection) *deliveryRepository {
	return &deliveryRepository{
		pool:      pool,
		txManager: txManager,
	}
}

// Create - создать доставку
func (dr *deliveryRepository) Create(ctx context.Context, del *delivery.AssignCreate) (*delivery.Delivery, error) {
	delCreateRow := domainToRowCreate(del)
	var delRow deliveryRow

	tx, err := dr.txManager.GetTx(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	query := `
	INSERT INTO delivery (courier_id, order_id, deadline)
	VALUES ($1, $2, $3)
	RETURNING id, courier_id, order_id, assigned_at, deadline;`
	args := []any{delCreateRow.CourierID, delCreateRow.OrderID, delCreateRow.Deadline}

	err = tx.QueryRow(
		ctx, query, args...,
	).Scan(&delRow.DeliveryID, &delRow.CourierID, &delRow.OrderID, &delRow.AssignedAt, &delRow.Deadline)

	if err != nil {
		return nil, mapError(err)
	}
	return rowToDomainDelivery(&delRow), nil
}

// Delete - удалить доставку
func (dr *deliveryRepository) Delete(ctx context.Context, orderID order.OrderID) (*delivery.Delivery, error) {
	var delRow deliveryRow

	tx, err := dr.txManager.GetTx(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	query := `
	DELETE FROM delivery WHERE order_id = $1
	RETURNING id, courier_id, order_id, assigned_at, deadline;`
	args := []any{orderID.OrderID}

	err = tx.QueryRow(
		ctx, query, args...,
	).Scan(&delRow.DeliveryID, &delRow.CourierID, &delRow.OrderID, &delRow.AssignedAt, &delRow.Deadline)

	if err != nil {
		return nil, mapError(err)
	}
	return rowToDomainDelivery(&delRow), nil
}

// Get - получить доставку
func (dr *deliveryRepository) Get(ctx context.Context, orderID order.OrderID) (*delivery.Delivery, error) {
	var delRow deliveryRow

	tx, err := dr.txManager.GetTx(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	query := `
	SELECT id, courier_id, order_id, assigned_at, deadline
	FROM delivery
	WHERE order_id = $1;`
	args := []any{orderID.OrderID}

	err = tx.QueryRow(
		ctx, query, args...,
	).Scan(&delRow.DeliveryID, &delRow.CourierID, &delRow.OrderID, &delRow.AssignedAt, &delRow.Deadline)

	if err != nil {
		return nil, mapError(err)
	}
	return rowToDomainDelivery(&delRow), nil
}
