package deliveryhttp

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/shft1/service-courier/internal/domain/delivery"
	"github.com/shft1/service-courier/observability/logger"
)

// DeliveryHandler - обработчик доставок
type DeliveryHandler struct {
	log     logger.Logger
	service deliveryService
}

// NewDeliveryHandler - конструктор обработчика доставок
func NewDeliveryHandler(log logger.Logger, service deliveryService) *DeliveryHandler {
	return &DeliveryHandler{
		service: service,
	}
}

// Assign - обрабатывает создание доставки на свободного курьера
func (dh *DeliveryHandler) Assign(w http.ResponseWriter, r *http.Request) {
	var orderReq DeliveryOrderRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		writeResponse(dh.log, w, 0, nil, delivery.ErrDeliveryInvalidData)
		return
	}
	if err := dh.validateOrderID(orderReq); err != nil {
		writeResponse(dh.log, w, 0, nil, err)
		return
	}
	del, err := dh.service.Assign(ctx, toDomainOrderID(orderReq))
	if err != nil {
		writeResponse(dh.log, w, 0, nil, err)
		return
	}
	writeResponse(dh.log, w, http.StatusOK, domainToDTOAssign(del), nil)
}

// Unassign - обрабатывает удаление доставки и освобождение соответствующего курьера
func (dh *DeliveryHandler) Unassign(w http.ResponseWriter, r *http.Request) {
	var orderID DeliveryOrderRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&orderID); err != nil {
		writeResponse(dh.log, w, 0, nil, delivery.ErrDeliveryInvalidData)
		return
	}
	if err := dh.validateOrderID(orderID); err != nil {
		writeResponse(dh.log, w, 0, nil, err)
		return
	}
	del, err := dh.service.Unassign(ctx, toDomainOrderID(orderID))
	if err != nil {
		writeResponse(dh.log, w, 0, nil, err)
		return
	}
	writeResponse(dh.log, w, http.StatusOK, domainToDTOUnassign(del), nil)
}

// validateOrderID - валидирует модель запроса
func (dh *DeliveryHandler) validateOrderID(orderID DeliveryOrderRequest) error {
	if orderID.OrderID == "" {
		return delivery.ErrDeliveryEmptyData
	}
	if _, err := uuid.Parse(orderID.OrderID); err != nil {
		return delivery.ErrDeliveryInvalidOrderID
	}
	return nil
}
