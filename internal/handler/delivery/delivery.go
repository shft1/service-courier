package delivery

import (
	"encoding/json"
	"net/http"
	"service-courier/internal/entity/delivery"

	"github.com/google/uuid"
)

type DeliveryHandler struct {
	service deliveryService
}

func NewDeliveryHandler(service deliveryService) *DeliveryHandler {
	return &DeliveryHandler{
		service: service,
	}
}

func (dh *DeliveryHandler) DeliveryAssign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var orderID delivery.DeliveryOrderID
	if err := json.NewDecoder(r.Body).Decode(&orderID); err != nil {
		deliveryMapResponse(w, 0, nil, delivery.ErrDeliveryInvalidData)
		return
	}
	if err := dh.validateOrderID(&orderID); err != nil {
		deliveryMapResponse(w, 0, nil, err)
		return
	}
	delivery, err := dh.service.DeliveryAssign(ctx, &orderID)
	deliveryMapResponse(w, http.StatusOK, delivery, err)
}

func (dh *DeliveryHandler) DeliveryUnassign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var orderID delivery.DeliveryOrderID
	if err := json.NewDecoder(r.Body).Decode(&orderID); err != nil {
		deliveryMapResponse(w, 0, nil, delivery.ErrDeliveryInvalidData)
		return
	}
	if err := dh.validateOrderID(&orderID); err != nil {
		deliveryMapResponse(w, 0, nil, err)
		return
	}
	delivery, err := dh.service.DeliveryUnassign(ctx, &orderID)
	deliveryMapResponse(w, http.StatusOK, delivery, err)
}

func (dh *DeliveryHandler) validateOrderID(orderID *delivery.DeliveryOrderID) error {
	if orderID.OrderID == "" {
		return delivery.ErrDeliveryEmptyData
	}
	if _, err := uuid.Parse(orderID.OrderID); err != nil {
		return delivery.ErrDeliveryInvalidOrderID
	}
	return nil
}
