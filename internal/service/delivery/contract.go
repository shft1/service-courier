package delivery

import (
	"context"
	"service-courier/internal/entity/courier"
	"service-courier/internal/entity/delivery"
)

//go:generate mockgen -source=contract.go -destination=mocks_test.go -package=delivery_test

type TxManagerDo interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type deliveryRepository interface {
	CreateDelivery(ctx context.Context, d *delivery.DeliveryCreate) (*delivery.DeliveryGet, error)
	DeleteDelivery(ctx context.Context, orderID *delivery.DeliveryOrderID) (*delivery.DeliveryGet, error)
	RecheckDelivery(ctx context.Context) error
}

type courierRepository interface {
	GetAvailable(ctx context.Context) (*courier.CourierGet, error)
	SetBusy(ctx context.Context, courierID int) error
	SetAvailable(ctx context.Context, courierID int) error
}
