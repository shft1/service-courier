package courier

type CourierCreate struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
}

type CourierGet struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
}

type CourierUpdate struct {
	ID     int     `json:"id"`
	Name   *string `json:"name"`
	Phone  *string `json:"phone"`
	Status *string `json:"status"`
}
