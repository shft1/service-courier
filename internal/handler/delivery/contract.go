package delivery

import (
	"context"
	"service-courier/internal/entity/delivery"
)

//go:generate mockgen -destination=./mocks/mock_delivery_service.go -package=mocks . deliveryService

type deliveryService interface {
	DeliveryAssign(ctx context.Context, orderID *delivery.DeliveryOrderID) (*delivery.DeliveryAssign, error)
	DeliveryUnassign(ctx context.Context, orderID *delivery.DeliveryOrderID) (*delivery.DeliveryUnassign, error)
	DeliveryCheck(ctx context.Context) error
}
