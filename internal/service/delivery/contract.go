package delivery

import (
	"context"
	"service-courier/internal/entity/courier"
	"service-courier/internal/entity/delivery"
)

//go:generate mockgen -destination=./mocks/mock_txmanager_do.go -package=mocks . TxManagerDo

type TxManagerDo interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

//go:generate mockgen -destination=./mocks/mock_delivery_repository.go -package=mocks . deliveryRepository

type deliveryRepository interface {
	CreateDelivery(ctx context.Context, d *delivery.DeliveryCreate) (*delivery.DeliveryGet, error)
	DeleteDelivery(ctx context.Context, orderID *delivery.DeliveryOrderID) (*delivery.DeliveryGet, error)
	RecheckDelivery(ctx context.Context) error
}

//go:generate mockgen -destination=./mocks/mock_courier_repository.go -package=mocks . courierRepository

type courierRepository interface {
	GetAvailable(ctx context.Context) (*courier.CourierGet, error)
	SetBusy(ctx context.Context, courierID int) error
	SetAvailable(ctx context.Context, courierID int) error
}
