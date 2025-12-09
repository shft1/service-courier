package deliveryhttp

import "time"

// DeliveryAssignResponse - модель возврата на создание доставки
type DeliveryAssignResponse struct {
	CourierID     int64     `json:"courier_id"`
	OrderID       string    `json:"order_id"`
	TransportType string    `json:"transport_type"`
	Deadline      time.Time `json:"delivery_deadline"`
}

// DeliveryUnassignResponse - модель возврата на удаление доставки
type DeliveryUnassignResponse struct {
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
	CourierID int64  `json:"courier_id"`
}

// DeliveryOrderRequest - модель запроса на создание доставки
type DeliveryOrderRequest struct {
	OrderID string `json:"order_id"`
}
