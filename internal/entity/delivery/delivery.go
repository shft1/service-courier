package delivery

import "time"

type DeliveryOrderID struct {
	OrderID string `json:"order_id"`
}

type DeliveryAssign struct {
	CourierID        int       `json:"courier_id"`
	OrderID          string    `json:"order_id"`
	TransportType    string    `json:"transport_type"`
	DeliveryDeadline time.Time `json:"delivery_deadline"`
}

type DeliveryUnassign struct {
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
	CourierID int    `json:"courier_id"`
}

type DeliveryGet struct {
	DeliveryID       int       `json:"delivery_id"`
	CourierID        int       `json:"courier_id"`
	OrderID          string    `json:"order_id"`
	AssignedAt       time.Time `json:"delivery_assign"`
	DeliveryDeadline time.Time `json:"delivery_deadline"`
}

type DeliveryCreate struct {
	CourierID        int       `json:"courier_id"`
	OrderID          string    `json:"order_id"`
	DeliveryDeadline time.Time `json:"delivery_deadline"`
}
