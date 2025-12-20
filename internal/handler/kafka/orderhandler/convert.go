package orderhandler

import "service-courier/internal/domain/order"

func dtoToDomainOrderID(dto *message) order.OrderID {
	return order.OrderID{OrderID: dto.OrderID}
}
