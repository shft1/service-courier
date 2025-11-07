package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"service-courier/internal/entity/courier"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type courierHandler struct {
	service courier.CourierService
}

func NewCourierHandler(service courier.CourierService) *courierHandler {
	return &courierHandler{
		service: service,
	}
}

func (ch *courierHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c courier.CourierCreate
	json.NewDecoder(r.Body).Decode(&c)
	if c.Name == "" || c.Phone == "" || c.Status == "" {
		http.Error(w, `{"error": "Missing required fields"}`, http.StatusBadRequest)
		return
	}
	if !regexp.MustCompile(`^\+?\d{10,16}$`).MatchString(c.Phone) {
		http.Error(w, `{"error": "Invalid Phone"}`, http.StatusBadRequest)
		return
	}
	err := ch.service.Create(&c)
	if err != nil {
		switch {
		case errors.Is(err, courier.ErrCourierExistPhone):
			http.Error(w, `{"error": "Courier with same phone already exists"}`, http.StatusConflict)
		default:
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (ch *courierHandler) Update(w http.ResponseWriter, r *http.Request) {
	var updCourier courier.CourierUpdate
	json.NewDecoder(r.Body).Decode(&updCourier)
	if updCourier.ID == nil {
		http.Error(w, `{"error": "The user's ID wasn't transmitted"}`, http.StatusBadRequest)
		return
	}
	if updCourier.Phone != nil {
		if !regexp.MustCompile(`^\+?\d{10,16}$`).MatchString(*updCourier.Phone) {
			http.Error(w, `{"error": "Invalid Phone"}`, http.StatusBadRequest)
			return
		}
	}
	err := ch.service.Update(&updCourier)
	if err != nil {
		switch {
		case errors.Is(err, courier.ErrCourierExistPhone):
			http.Error(w, `{"error": "Courier with same phone already exists"}`, http.StatusConflict)
		case errors.Is(err, courier.ErrCourierNotFound):
			http.Error(w, `{"error": "Courier not found"}`, http.StatusNotFound)
		default:
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (ch *courierHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid courier ID"}`, http.StatusBadRequest)
		return
	}
	c, err := ch.service.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, courier.ErrCourierNotFound):
			http.Error(w, `{"error": "Courier not found"}`, http.StatusNotFound)
		case errors.Is(err, courier.ErrCourierInvalidID):
			http.Error(w, `{"error": "Invalid courier ID"}`, http.StatusBadRequest)
		default:
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func (ch *courierHandler) GetMulti(w http.ResponseWriter, r *http.Request) {
	couriers, err := ch.service.GetMulti()
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(couriers)
}
