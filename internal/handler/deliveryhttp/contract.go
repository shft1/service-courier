package deliveryhttp

import (
	"context"
	"service-courier/internal/domain/delivery"
)

//go:generate mockgen -source=contract.go -destination=./mocks_test.go -package=deliveryhttp_test

type deliveryService interface {
	Assign(ctx context.Context, orderID delivery.OrderID) (*delivery.AssignResult, error)
	Unassign(ctx context.Context, orderID delivery.OrderID) (*delivery.UnassignResult, error)
	CheckDelivery(ctx context.Context) error
}
