package courier

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"service-courier/internal/entity/courier"
)

type courierRepository interface {
	Create(ctx context.Context, c *courier.CourierCreate) error
	GetByID(ctx context.Context, id int) (*courier.CourierGet, error)
	GetMulti(ctx context.Context) ([]courier.CourierGet, error)
	Update(ctx context.Context, c *courier.CourierUpdate) error
}

type courierService struct {
	repository courierRepository
}

func NewCourierService(repo courierRepository) *courierService {
	return &courierService{
		repository: repo,
	}
}

func (cs *courierService) Create(ctx context.Context, c *courier.CourierCreate) error {
	if err := cs.validateCreate(c); err != nil {
		return err
	}

	err := cs.repository.Create(ctx, c)
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

func (cs *courierService) Update(ctx context.Context, c *courier.CourierUpdate) error {
	if err := cs.validateUpdate(c); err != nil {
		return err
	}

	err := cs.repository.Update(ctx, c)
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

func (cs *courierService) GetByID(ctx context.Context, id int) (*courier.CourierGet, error) {
	if id < 1 {
		return nil, courier.ErrCourierInvalidID
	}

	c, err := cs.repository.GetByID(ctx, id)
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

func (cs *courierService) GetMulti(ctx context.Context) ([]courier.CourierGet, error) {
	couriers, err := cs.repository.GetMulti(ctx)

	if err != nil {
		return nil, fmt.Errorf("service: failed to get couriers: %w", err)
	}
	return couriers, nil
}

func (cs *courierService) validateCreate(c *courier.CourierCreate) error {
	if c.Name == "" || c.Phone == "" || c.Status == "" {
		return courier.ErrCourierEmptyData
	}
	if !regexp.MustCompile(`^\+?\d{10,16}$`).MatchString(c.Phone) {
		return courier.ErrCourierInvalidPhone
	}
	return nil
}

func (cs *courierService) validateUpdate(c *courier.CourierUpdate) error {
	if c.ID < 1 {
		return courier.ErrCourierInvalidID
	}
	if c.Phone != nil && !regexp.MustCompile(`^\+?\d{10,16}$`).MatchString(*c.Phone) {
		return courier.ErrCourierInvalidPhone
	}
	return nil
}
