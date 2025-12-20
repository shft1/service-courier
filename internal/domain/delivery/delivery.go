package delivery

import "time"

const UnassignStatus = "unassigned"

// Delivery - доменная сущность
type Delivery struct {
	DeliveryID int64
	CourierID  int64
	OrderID    string
	AssignedAt time.Time
	Deadline   time.Time
}

// AssignResult - объект успешного создания доставки
type AssignResult struct {
	CourierID     int64
	OrderID       string
	TransportType string
	Deadline      time.Time
}

// UnassignResult - объект успешного удаления доставки
type UnassignResult struct {
	CourierID int64
	OrderID   string
	Status    string
}

// AssignCreate - объект создания доставки
type AssignCreate struct {
	CourierID int64
	OrderID   string
	Deadline  time.Time
}

// CompleteResult - объект успешного выполнения доставки
type CompleteResult struct {
	CourierID int64
	OrderID   string
	Deadline  time.Time
}
