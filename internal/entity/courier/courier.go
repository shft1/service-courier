package courier

type CourierCreate struct {
	Name          string  `json:"name"`
	Phone         string  `json:"phone"`
	Status        *string `json:"status"`
	TransportType *string `json:"transport_type"`
}

type CourierGet struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Status        string `json:"status"`
	TransportType string `json:"transport_type"`
}

type CourierUpdate struct {
	ID            int     `json:"id"`
	Name          *string `json:"name"`
	Phone         *string `json:"phone"`
	Status        *string `json:"status"`
	TransportType *string `json:"transport_type"`
}
