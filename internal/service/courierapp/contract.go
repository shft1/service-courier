package courierapp

import (
	"context"
	"service-courier/internal/domain/courier"
)

//go:generate mockgen -source=contract.go -destination=./mocks_test.go -package=courierapp_test

type courierRepository interface {
	Create(ctx context.Context, cour *courier.CourierCreate) (int64, error)
	Update(ctx context.Context, cour *courier.CourierUpdate) (int64, error)
	GetByID(ctx context.Context, id int64) (*courier.Courier, error)
	GetMulti(ctx context.Context) ([]courier.Courier, error)
	GetAvailable(ctx context.Context) (*courier.Courier, error)
	SetBusy(ctx context.Context, id int64) (int64, error)
	SetAvailable(ctx context.Context, courierID int64) (int64, error)
	ReleaseStaleBusy(ctx context.Context) error
}
