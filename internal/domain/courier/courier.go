package courier

import "time"

// Courier - доменная сущность
type Courier struct {
	ID            int64
	Name          string
	Phone         string
	Status        string
	TransportType string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// CourierCreate - бизнес объект создания
type CourierCreate struct {
	Name          string
	Phone         string
	Status        *string
	TransportType *string
}

// CourierUpdate - бизнес объект обновления
type CourierUpdate struct {
	ID            int64
	Name          *string
	Phone         *string
	Status        *string
	TransportType *string
}
