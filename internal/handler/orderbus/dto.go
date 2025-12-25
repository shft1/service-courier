package orderbus

import "time"

// message - модель события из Kafka
type message struct {
	OrderID   string    `json:"order_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
