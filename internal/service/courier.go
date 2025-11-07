package service

import (
	"errors"
	"fmt"
	"service-courier/internal/entity/courier"
)

type courierService struct {
	repository courier.CourierRepository
}

func NewCourierService(repo courier.CourierRepository) courier.CourierService {
	return &courierService{
		repository: repo,
	}
}

func (cs *courierService) Create(c *courier.CourierCreate) error {
	err := cs.repository.Create(c)
	if err != nil {
		switch {
		case errors.Is(err, courier.ErrCourierExistPhone):
			return courier.ErrCourierExistPhone
		default:
			return fmt.Errorf("service: failed to create courier: %w", err)
		}
	}
	return nil
}

func (cs *courierService) Update(c *courier.CourierUpdate) error {
	err := cs.repository.Update(c)
	if err != nil {
		switch {
		case errors.Is(err, courier.ErrCourierExistPhone):
			return courier.ErrCourierExistPhone
		case errors.Is(err, courier.ErrCourierNotFound):
			return courier.ErrCourierNotFound
		default:
			return fmt.Errorf("service: failed to update courier: %w", err)
		}
	}
	return nil
}

func (cs *courierService) GetByID(id int) (*courier.CourierGet, error) {
	if id < 1 {
		return nil, courier.ErrCourierInvalidID
	}
	c, err := cs.repository.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, courier.ErrCourierNotFound):
			return nil, courier.ErrCourierNotFound
		default:
			return nil, fmt.Errorf("service: failed to get courier by id: %w", err)
		}
	}
	return c, nil
}

func (cs *courierService) GetMulti() ([]courier.CourierGet, error) {
	couriers, err := cs.repository.GetMulti()
	if err != nil {
		return nil, fmt.Errorf("service: failed to get couriers: %w", err)
	}
	return couriers, nil
}
