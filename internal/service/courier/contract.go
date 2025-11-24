package courier

import (
	"context"
	"service-courier/internal/entity/courier"
)

type сourierRepository interface {
	Create(ctx context.Context, c *courier.CourierCreate) error
	GetByID(ctx context.Context, id int) (*courier.CourierGet, error)
	GetMulti(ctx context.Context) ([]courier.CourierGet, error)
	Update(ctx context.Context, c *courier.CourierUpdate) error
	GetAvailable(ctx context.Context) (*courier.CourierGet, error)
	SetBusy(ctx context.Context, id int) error
	SetAvailable(ctx context.Context, courierID int) error
}
