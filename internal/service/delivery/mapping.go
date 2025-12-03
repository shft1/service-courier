package delivery

import (
	"errors"
	"fmt"
	"service-courier/internal/entity/courier"
	"service-courier/internal/entity/delivery"
)

func deliveryServiceMapError(err error) error {
	switch {
	case errors.Is(err, delivery.ErrDeliveryExist):
		return delivery.ErrDeliveryExist
	case errors.Is(err, courier.ErrCourierNotFound):
		return delivery.ErrDeliveryCourierLost
	case errors.Is(err, delivery.ErrDeliveryInvalidAssignCourier):
		return delivery.ErrDeliveryInvalidAssignCourier
	case errors.Is(err, courier.ErrCourierAvailable):
		return delivery.ErrDeliveryNotAvailableCourier
	case errors.Is(err, delivery.ErrDeliveryNotFound):
		return delivery.ErrDeliveryNotFound
	default:
		return fmt.Errorf("service: failed to work with delivery: %w", err)
	}
}
