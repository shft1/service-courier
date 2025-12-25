package ordergrpc

import (
	"service-courier/internal/domain/order"
	"service-courier/internal/proto/orderpb"
)

// toDomainOrderList - конвертатор в список доменных сущностей заказов
func toDomainOrderList(ordersDTO []*orderpb.Order) []*order.Order {
	ordersResp := make([]*order.Order, 0, len(ordersDTO))
	for _, orderDTO := range ordersDTO {
		ordersResp = append(ordersResp, toDomainOrder(orderDTO))
	}
	return ordersResp
}

// toDomainOrder - конвертатор в доменную сущность заказа
func toDomainOrder(orderDTO *orderpb.Order) *order.Order {
	return &order.Order{
		OrderID:   orderDTO.Id,
		Status:    orderDTO.Status,
		CreatedAt: orderDTO.CreatedAt.AsTime(),
	}
}
