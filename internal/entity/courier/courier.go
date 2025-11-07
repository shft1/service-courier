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
	ID     *int    `json:"id"`
	Name   *string `json:"name"`
	Phone  *string `json:"phone"`
	Status *string `json:"status"`
}

type CourierService interface {
	Create(c *CourierCreate) error
	GetByID(id int) (*CourierGet, error)
	GetMulti() ([]CourierGet, error)
	Update(c *CourierUpdate) error
}

type CourierRepository interface {
	Create(c *CourierCreate) error
	GetByID(id int) (*CourierGet, error)
	GetMulti() ([]CourierGet, error)
	Update(c *CourierUpdate) error
}
