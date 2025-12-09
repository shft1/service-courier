package courierhttp

import (
	"encoding/json"
	"net/http"
	"regexp"
	"service-courier/internal/domain/courier"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// CourierHandler - обработчик курьеров
type CourierHandler struct {
	service courierService
}

// NewCourierHandler - конструктор обработчика курьеров
func NewCourierHandler(service courierService) *CourierHandler {
	return &CourierHandler{
		service: service,
	}
}

// Create - обрабатывает запрос на создание курьера
func (ch *CourierHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var courReq CourierCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&courReq); err != nil {
		writeResponse(w, 0, nil, courier.ErrCourierInvalidData)
		return
	}
	if err := ch.validateCreate(&courReq); err != nil {
		writeResponse(w, 0, nil, err)
		return
	}
	_, err := ch.service.Create(ctx, toDomainCreate(&courReq))
	writeResponse(w, http.StatusCreated, nil, err)
}

// Update - обрабатывает запрос на частичное обновление курьера
func (ch *CourierHandler) Update(w http.ResponseWriter, r *http.Request) {
	var courReq CourierUpdateRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&courReq); err != nil {
		writeResponse(w, 0, nil, courier.ErrCourierInvalidData)
		return
	}
	if err := ch.validateUpdate(&courReq); err != nil {
		writeResponse(w, 0, nil, err)
		return
	}
	_, err := ch.service.Update(ctx, toDomainUpdate(&courReq))
	writeResponse(w, http.StatusOK, nil, err)
}

// GetByID - обрабатывает запрос на получение курьера по его ID
func (ch *CourierHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		writeResponse(w, 0, nil, courier.ErrCourierInvalidID)
		return
	}
	cour, err := ch.service.GetByID(ctx, int64(id))
	if err != nil {
		writeResponse(w, 0, nil, err)
		return
	}
	writeResponse(w, http.StatusOK, domainToDTO(cour), nil)
}

// GetMulti - обрабатывает запрос на получение курьеров
func (ch *CourierHandler) GetMulti(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cours, err := ch.service.GetMulti(ctx)
	if err != nil {
		writeResponse(w, 0, nil, err)
		return
	}
	writeResponse(w, http.StatusOK, domainToDTOList(cours), nil)
}

// validateCreate - валидирует модель создания курьера
func (ch *CourierHandler) validateCreate(courReq *CourierCreateRequest) error {
	if courReq.Name == "" || courReq.Phone == "" {
		return courier.ErrCourierEmptyData
	}
	if !regexp.MustCompile(`^\+?\d{10,16}$`).MatchString(courReq.Phone) {
		return courier.ErrCourierInvalidPhone
	}
	return nil
}

// validateUpdate - валидирует модель обновления курьера
func (ch *CourierHandler) validateUpdate(courReq *CourierUpdateRequest) error {
	if courReq.ID < 1 {
		return courier.ErrCourierInvalidID
	}
	if courReq.Phone != nil && !regexp.MustCompile(`^\+?\d{10,16}$`).MatchString(*courReq.Phone) {
		return courier.ErrCourierInvalidPhone
	}
	return nil
}
