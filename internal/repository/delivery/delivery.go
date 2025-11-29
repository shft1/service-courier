package delivery

import (
	"context"
	"service-courier/internal/entity/delivery"

	"github.com/jackc/pgx/v5/pgxpool"
)

type deliveryRepository struct {
	pool      *pgxpool.Pool
	txManager txManagerGetConnection
}

func NewDeliveryRepository(pool *pgxpool.Pool, txManager txManagerGetConnection) *deliveryRepository {
	return &deliveryRepository{
		pool:      pool,
		txManager: txManager,
	}
}

func (dr *deliveryRepository) CreateDelivery(ctx context.Context, d *delivery.DeliveryCreate) (*delivery.DeliveryGet, error) {
	var delivery delivery.DeliveryGet

	tx, err := dr.txManager.GetTx(ctx)
	if err != nil {
		return nil, deliveryRepoMapError(err)
	}
	query := `
	INSERT INTO delivery (courier_id, order_id, deadline)
	VALUES ($1, $2, $3)
	RETURNING id, courier_id, order_id, assigned_at, deadline;`
	args := []any{d.CourierID, d.OrderID, d.DeliveryDeadline}

	err = tx.QueryRow(
		ctx, query, args...,
	).Scan(&delivery.DeliveryID, &delivery.CourierID, &delivery.OrderID, &delivery.AssignedAt, &delivery.DeliveryDeadline)

	if err != nil {
		return nil, deliveryRepoMapError(err)
	}
	return &delivery, nil
}

func (dr *deliveryRepository) DeleteDelivery(ctx context.Context, orderID *delivery.DeliveryOrderID) (*delivery.DeliveryGet, error) {
	var delivery delivery.DeliveryGet

	tx, err := dr.txManager.GetTx(ctx)
	if err != nil {
		return nil, deliveryRepoMapError(err)
	}
	query := `
	DELETE FROM delivery WHERE order_id = $1
	RETURNING id, courier_id, order_id, assigned_at, deadline;`
	args := []any{orderID.OrderID}

	err = tx.QueryRow(
		ctx, query, args...,
	).Scan(&delivery.DeliveryID, &delivery.CourierID, &delivery.OrderID, &delivery.AssignedAt, &delivery.DeliveryDeadline)

	if err != nil {
		return nil, deliveryRepoMapError(err)
	}
	return &delivery, nil
}

func (dr *deliveryRepository) RecheckDelivery(ctx context.Context) error {
	query := `
	UPDATE couriers
	SET status = 'available'
	WHERE status = 'busy' AND id NOT IN (
		SELECT courier_id
		FROM delivery
		WHERE NOW() <= deadline
	);`
	if _, err := dr.pool.Exec(ctx, query); err != nil {
		return deliveryRepoMapError(err)
	}
	return nil
}
