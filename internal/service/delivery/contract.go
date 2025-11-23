package delivery

import (
	"context"
	"service-courier/internal/entity/courier"
	"service-courier/internal/entity/delivery"
)

type deliveryRepository interface {
	CreateDelivery(ctx context.Context, d *delivery.DeliveryCreate) (*delivery.DeliveryGet, error)
	DeleteDelivery(ctx context.Context, orderID *delivery.DeliveryOrderID) (*delivery.DeliveryGet, error)
}

type courierRepository interface {
	GetAvailable(ctx context.Context) (*courier.CourierGet, error)
	SetBusy(ctx context.Context, courierID int) error
	SetAvailable(ctx context.Context, courierID int) error
}
