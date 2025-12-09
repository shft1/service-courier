package deliveryapp

import (
	"context"
	"service-courier/internal/domain/courier"
	"service-courier/internal/domain/delivery"
	"time"
)

//go:generate mockgen -source=contract.go -destination=mocks_test.go -package=deliveryapp_test

type deliveryRepository interface {
	Create(ctx context.Context, del *delivery.AssignCreate) (*delivery.Delivery, error)
	Delete(ctx context.Context, orderID delivery.OrderID) (*delivery.Delivery, error)
}

type courierRepository interface {
	GetAvailable(ctx context.Context) (*courier.Courier, error)
	SetBusy(ctx context.Context, courierID int64) (int64, error)
	SetAvailable(ctx context.Context, courierID int64) (int64, error)
	ReleaseStaleBusy(ctx context.Context) error
}

type txManagerDo interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type timeCalculatorFactory interface {
	GetDeliveryCalculator(transportType string) (TimeCalculator, error)
}

type TimeCalculator interface {
	Calculate() time.Time
}
