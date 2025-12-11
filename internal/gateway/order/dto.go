package order

import (
	"service-courier/internal/proto/order"
	"time"
)

// OrderResponse - модель ответа от client pb
type OrderResponse struct {
	OrderID   string
	CreatedAt time.Time
}

// toOrderResponseList - конвертатор dto pb -> dto для воркера
func toOrderResponseList(ordersDTO []*order.Order) []*OrderResponse {
	ordersResp := make([]*OrderResponse, 0, len(ordersDTO))
	for _, orderDTO := range ordersDTO {
		ordersResp = append(ordersResp, &OrderResponse{
			OrderID:   orderDTO.Id,
			CreatedAt: orderDTO.CreatedAt.AsTime(),
		})
	}
	return ordersResp
}
