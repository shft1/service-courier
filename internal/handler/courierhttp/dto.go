package courierhttp

// CourierResponse - модель возврата курьера
type CourierResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Status        string `json:"status"`
	TransportType string `json:"transport_type"`
}

// CourierCreateRequest - модель запроса на создание курьера
type CourierCreateRequest struct {
	Name          string  `json:"name"`
	Phone         string  `json:"phone"`
	Status        *string `json:"status"`
	TransportType *string `json:"transport_type"`
}

// CourierUpdateRequest - модель запроса на обновление курьера
type CourierUpdateRequest struct {
	ID            int64   `json:"id"`
	Name          *string `json:"name"`
	Phone         *string `json:"phone"`
	Status        *string `json:"status"`
	TransportType *string `json:"transport_type"`
}
