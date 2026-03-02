package courierhttp

import (
	"context"

	"github.com/shft1/service-courier/internal/domain/courier"
)

//go:generate mockgen -source=contract.go -destination=./mocks_test.go -package=courierhttp_test

type courierService interface {
	Create(ctx context.Context, c *courier.CourierCreate) (int64, error)
	Update(ctx context.Context, c *courier.CourierUpdate) (int64, error)
	GetByID(ctx context.Context, id int64) (*courier.Courier, error)
	GetMulti(ctx context.Context) ([]courier.Courier, error)
}
