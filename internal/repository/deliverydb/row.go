package deliverydb

import "time"

// deliveryRow - модель доставки БД
type deliveryRow struct {
	DeliveryID int64     `db:"id"`
	CourierID  int64     `db:"courier_id"`
	OrderID    string    `db:"order_id"`
	AssignedAt time.Time `db:"assigned_at"`
	Deadline   time.Time `db:"deadline"`
}

// deliveryCreateRow - модель создания доставки в БД
type deliveryCreateRow struct {
	CourierID int64     `db:"courier_id"`
	OrderID   string    `db:"order_id"`
	Deadline  time.Time `db:"deadline"`
}
