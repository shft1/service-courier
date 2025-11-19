package courier

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"service-courier/internal/entity/courier"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func courierHandMapError(err error) (int, string) {
	var (
		status int
		msg    string
	)
	switch err {
	case courier.ErrCourierInvalidData:
		status = http.StatusBadRequest
		msg = err.Error()
	case courier.ErrCourierExistPhone:
		status = http.StatusConflict
		msg = err.Error()
	case courier.ErrCourierEmptyData:
		status = http.StatusBadRequest
		msg = err.Error()
	case courier.ErrCourierInvalidPhone:
		status = http.StatusBadRequest
		msg = err.Error()
	case courier.ErrCourierNotFound:
		status = http.StatusNotFound
		msg = err.Error()
	case courier.ErrCourierInvalidID:
		status = http.StatusBadRequest
		msg = err.Error()
	default:
		status = http.StatusInternalServerError
		msg = courier.ErrDatabase.Error()
	}
	return status, msg
}

func courierMapResponse(w http.ResponseWriter, status int, data any, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		stErr, msg := courierHandMapError(err)
		w.WriteHeader(stErr)
		json.NewEncoder(w).Encode(map[string]string{"error": msg})
		return
	}
	if data == nil {
		w.WriteHeader(status)
		return
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

type courierService interface {
	Create(ctx context.Context, c *courier.CourierCreate) error
	GetByID(ctx context.Context, id int) (*courier.CourierGet, error)
	GetMulti(ctx context.Context) ([]courier.CourierGet, error)
	Update(ctx context.Context, c *courier.CourierUpdate) error
}

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
	if err != nil {
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
	if c.Name == "" || c.Phone == "" || c.Status == "" {
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
