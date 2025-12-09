package courierdb

import "time"

// courierRow - модель курьера БД
type courierRow struct {
	ID            int64     `db:"id"`
	Name          string    `db:"name"`
	Phone         string    `db:"phone"`
	Status        string    `db:"status"`
	TransportType string    `db:"transport_type"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// courierCreateRow - модель создания курьера в БД
type courierCreateRow struct {
	Name          string  `db:"name"`
	Phone         string  `db:"phone"`
	Status        *string `db:"status"`
	TransportType *string `db:"transport_type"`
}

// courierUpdateRow - модель обновления курьера в БД
type courierUpdateRow struct {
	ID            int64   `db:"id"`
	Name          *string `db:"name"`
	Phone         *string `db:"phone"`
	Status        *string `db:"status"`
	TransportType *string `db:"transport_type"`
}
