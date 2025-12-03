package delivery

import (
	"context"
	"service-courier/internal/entity/delivery"
)

//go:generate mockgen -source=contract.go -destination=./mocks_test.go -package=delivery_test

type deliveryService interface {
	DeliveryAssign(ctx context.Context, orderID *delivery.DeliveryOrderID) (*delivery.DeliveryAssign, error)
	DeliveryUnassign(ctx context.Context, orderID *delivery.DeliveryOrderID) (*delivery.DeliveryUnassign, error)
	DeliveryCheck(ctx context.Context) error
}
