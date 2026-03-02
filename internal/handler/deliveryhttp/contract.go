package deliveryhttp

import (
	"context"

	"github.com/shft1/service-courier/internal/domain/delivery"
	"github.com/shft1/service-courier/internal/domain/order"
)

//go:generate mockgen -source=contract.go -destination=./mocks_test.go -package=deliveryhttp_test

type deliveryService interface {
	Assign(ctx context.Context, orderID order.OrderID) (*delivery.AssignResult, error)
	Unassign(ctx context.Context, orderID order.OrderID) (*delivery.UnassignResult, error)
	CheckDelivery(ctx context.Context) error
}
