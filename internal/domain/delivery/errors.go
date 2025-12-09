package delivery

import "errors"

var (
	ErrDeliveryNotAvailableCourier  = errors.New("all couriers are busy")
	ErrDeliveryEmptyData            = errors.New("the order ID can't be empty")
	ErrDeliveryExist                = errors.New("delivery with such a order ID already exists")
	ErrDeliveryNotFound             = errors.New("delivery with such a order ID wasn't found")
	ErrDeliveryInvalidData          = errors.New("order information is incorrect")
	ErrDeliveryInvalidOrderID       = errors.New("order ID is incorrect")
	ErrDeliveryInvalidAssignCourier = errors.New("couldn't assign a free courier to the order")
	ErrDeliveryCourierLost          = errors.New("connection of the courier executing the delivery is lost")
	ErrDeliveryDatabase             = errors.New("database error")
)
