package delivery

import (
	"context"
	"service-courier/internal/entity/delivery"
	"service-courier/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type deliveryRepository struct {
	pool      *pgxpool.Pool
	txManager repository.TxManagerGetConnection
}

func NewDeliveryRepository(pool *pgxpool.Pool, txManager repository.TxManagerGetConnection) *deliveryRepository {
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
