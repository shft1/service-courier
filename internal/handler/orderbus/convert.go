package orderbus

import "github.com/shft1/service-courier/internal/domain/order"

func dtoToDomainOrderID(dto *message) order.OrderID {
	return order.OrderID{OrderID: dto.OrderID}
}
