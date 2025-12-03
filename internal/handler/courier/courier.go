package courier

import (
	"encoding/json"
	"net/http"
	"regexp"
	"service-courier/internal/entity/courier"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CourierHandler struct {
	service courierService
}

func NewCourierHandler(service courierService) *CourierHandler {
	return &CourierHandler{
		service: service,
	}
}

func (ch *CourierHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var c courier.CourierCreate
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		courierMapResponse(w, 0, nil, courier.ErrCourierInvalidData)
		return
	}
	if err := ch.validateCreate(&c); err != nil {
		courierMapResponse(w, 0, nil, err)
		return
	}
	err := ch.service.Create(ctx, &c)
	courierMapResponse(w, http.StatusCreated, nil, err)
}

func (ch *CourierHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var updCourier courier.CourierUpdate
	if err := json.NewDecoder(r.Body).Decode(&updCourier); err != nil {
		courierMapResponse(w, 0, nil, courier.ErrCourierInvalidData)
		return
	}
	if err := ch.validateUpdate(&updCourier); err != nil {
		courierMapResponse(w, 0, nil, err)
		return
	}
	err := ch.service.Update(ctx, &updCourier)
	courierMapResponse(w, http.StatusOK, nil, err)
}

func (ch *CourierHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		courierMapResponse(w, 0, nil, courier.ErrCourierInvalidID)
		return
	}
	c, err := ch.service.GetByID(ctx, id)
	courierMapResponse(w, http.StatusOK, c, err)
}

func (ch *CourierHandler) GetMulti(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	couriers, err := ch.service.GetMulti(ctx)
	courierMapResponse(w, http.StatusOK, couriers, err)
}

func (ch *CourierHandler) validateCreate(c *courier.CourierCreate) error {
	if c.Name == "" || c.Phone == "" {
		return courier.ErrCourierEmptyData
	}
	if !regexp.MustCompile(`^\+?\d{10,16}$`).MatchString(c.Phone) {
		return courier.ErrCourierInvalidPhone
	}
	return nil
}

func (ch *CourierHandler) validateUpdate(c *courier.CourierUpdate) error {
	if c.ID < 1 {
		return courier.ErrCourierInvalidID
	}
	if c.Phone != nil && !regexp.MustCompile(`^\+?\d{10,16}$`).MatchString(*c.Phone) {
		return courier.ErrCourierInvalidPhone
	}
	return nil
}
