package deliveryapp

import (
	"context"
	"time"

	"github.com/shft1/service-courier/internal/domain/courier"
	"github.com/shft1/service-courier/internal/domain/delivery"
	"github.com/shft1/service-courier/internal/domain/order"
)

//go:generate mockgen -source=contract.go -destination=mocks_test.go -package=deliveryapp_test

type deliveryRepository interface {
	Create(ctx context.Context, del *delivery.AssignCreate) (*delivery.Delivery, error)
	Delete(ctx context.Context, orderID order.OrderID) (*delivery.Delivery, error)
	Get(ctx context.Context, orderID order.OrderID) (*delivery.Delivery, error)
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

type Executor interface {
	Execute(ctx context.Context, orderID order.OrderID) error
}

type deliveryExecutor interface {
	Assign(ctx context.Context, orderID order.OrderID) (*delivery.AssignResult, error)
	Unassign(ctx context.Context, orderID order.OrderID) (*delivery.UnassignResult, error)
	Complete(ctx context.Context, orderID order.OrderID) (*delivery.CompleteResult, error)
}

type DeliveryAssign interface {
	Assign(ctx context.Context, orderID order.OrderID) (*delivery.AssignResult, error)
}

type DeliveryUnassign interface {
	Unassign(ctx context.Context, orderID order.OrderID) (*delivery.UnassignResult, error)
}

type DeliveryComplete interface {
	Complete(ctx context.Context, orderID order.OrderID) (*delivery.CompleteResult, error)
}
