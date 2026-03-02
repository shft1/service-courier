package courierapp

import (
	"context"

	"github.com/shft1/service-courier/internal/domain/courier"
)

// courierService - cервис курьера
type courierService struct {
	repository courierRepository
}

// NewCourierService - конструктор сервиса курьера
func NewCourierService(repo courierRepository) *courierService {
	return &courierService{
		repository: repo,
	}
}

// Create - создать курьера
func (cs *courierService) Create(ctx context.Context, cour *courier.CourierCreate) (int64, error) {
	id, err := cs.repository.Create(ctx, cour)

	if err != nil {
		return -1, mapError(err)
	}
	return id, nil
}

// Update - обновить курьера
func (cs *courierService) Update(ctx context.Context, cour *courier.CourierUpdate) (int64, error) {
	id, err := cs.repository.Update(ctx, cour)

	if err != nil {
		return -1, mapError(err)
	}
	return id, nil
}

// GetByID - получить курьера по ID
func (cs *courierService) GetByID(ctx context.Context, id int64) (*courier.Courier, error) {
	cour, err := cs.repository.GetByID(ctx, id)

	if err != nil {
		return nil, mapError(err)
	}
	return cour, nil
}

// GetMulti - получить курьеров
func (cs *courierService) GetMulti(ctx context.Context) ([]courier.Courier, error) {
	cours, err := cs.repository.GetMulti(ctx)

	if err != nil {
		return nil, mapError(err)
	}
	return cours, nil
}
