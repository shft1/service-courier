package courierapp

import (
	"errors"
	"fmt"

	"service-courier/internal/domain/courier"
)

// mapError - маппинг ошибок сервиса курьеров
func mapError(err error) error {
	switch {
	case errors.Is(err, courier.ErrCourierExistPhone):
		return courier.ErrCourierExistPhone
	case errors.Is(err, courier.ErrCourierNotFound):
		return courier.ErrCourierNotFound
	default:
		return fmt.Errorf("service: failed to work with courier: %w", err)
	}
}
