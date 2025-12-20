package orderhandler

import (
	"context"
	"service-courier/internal/domain/order"
	"service-courier/internal/service/deliveryapp"
)

type eventStrategyFactory interface {
	GetEventStrategy(statusMsg string, statusNow string) (deliveryapp.Executor, error)
}

type orderGateway interface {
	GetOrderByID(ctx context.Context, orderID order.OrderID) (*order.Order, error)
}
