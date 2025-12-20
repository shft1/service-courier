package order

import "time"

// Order - доменная сущность заказа
type Order struct {
	OrderID   string
	Status    string
	CreatedAt time.Time
}

// OrderID - модель номера заказа
type OrderID struct {
	OrderID string
}
