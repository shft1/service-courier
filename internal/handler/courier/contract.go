package courier

import (
	"context"
	"service-courier/internal/entity/courier"
)

//go:generate mockgen -destination=./mocks/mock_courier_service.go -package=mocks . courierService

type courierService interface {
	Create(ctx context.Context, c *courier.CourierCreate) error
	GetByID(ctx context.Context, id int) (*courier.CourierGet, error)
	GetMulti(ctx context.Context) ([]courier.CourierGet, error)
	Update(ctx context.Context, c *courier.CourierUpdate) error
}
