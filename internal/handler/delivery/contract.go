package delivery

import (
	"context"
	"service-courier/internal/entity/delivery"
)

type deliveryService interface {
	DeliveryAssign(ctx context.Context, orderID *delivery.DeliveryOrderID) (*delivery.DeliveryAssign, error)
	DeliveryUnassign(ctx context.Context, orderID *delivery.DeliveryOrderID) (*delivery.DeliveryUnassign, error)
	DeliveryCheck(ctx context.Context) error
}
