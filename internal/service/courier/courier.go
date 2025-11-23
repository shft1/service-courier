package courier

import (
	"context"
	"service-courier/internal/entity/courier"
)

type courierService struct {
	repository сourierRepository
}

func NewCourierService(repo сourierRepository) *courierService {
	return &courierService{
		repository: repo,
	}
}

func (cs *courierService) Create(ctx context.Context, c *courier.CourierCreate) error {
	err := cs.repository.Create(ctx, c)

	if err != nil {
		return courierServiceMapError(err)
	}
	return nil
}

func (cs *courierService) Update(ctx context.Context, c *courier.CourierUpdate) error {
	err := cs.repository.Update(ctx, c)

	if err != nil {
		return courierServiceMapError(err)
	}
	return nil
}

func (cs *courierService) GetByID(ctx context.Context, id int) (*courier.CourierGet, error) {
	c, err := cs.repository.GetByID(ctx, id)

	if err != nil {
		return nil, courierServiceMapError(err)
	}
	return c, nil
}

func (cs *courierService) GetMulti(ctx context.Context) ([]courier.CourierGet, error) {
	couriers, err := cs.repository.GetMulti(ctx)

	if err != nil {
		return nil, courierServiceMapError(err)
	}
	return couriers, nil
}
