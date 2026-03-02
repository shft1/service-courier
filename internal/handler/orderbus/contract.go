package orderbus

import (
	"context"

	"github.com/shft1/service-courier/internal/domain/order"
	"github.com/shft1/service-courier/internal/service/deliveryapp"
)

type eventStrategyFactory interface {
	GetEventStrategy(statusMsg string, statusNow string) (deliveryapp.Executor, error)
}

type orderGateway interface {
	GetOrderByID(ctx context.Context, orderID order.OrderID) (*order.Order, error)
}
