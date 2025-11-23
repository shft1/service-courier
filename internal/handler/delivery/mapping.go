package delivery

import (
	"encoding/json"
	"net/http"
	"service-courier/internal/entity/delivery"
)

func deliveryHandMapError(err error) (int, string) {
	var (
		status int
		msg    string
	)
	switch err {
	case delivery.ErrDeliveryExist:
		status = http.StatusConflict
		msg = err.Error()
	case delivery.ErrDeliveryNotFound:
		status = http.StatusNotFound
		msg = err.Error()
	case delivery.ErrDeliveryNotAvailableCourier:
		status = http.StatusConflict
		msg = err.Error()
	case delivery.ErrDeliveryInvalidAssignCourier:
		status = http.StatusConflict
		msg = err.Error()
	case delivery.ErrDeliveryInvalidOrderID:
		status = http.StatusBadRequest
		msg = err.Error()
	case delivery.ErrDeliveryEmptyData:
		status = http.StatusBadRequest
		msg = err.Error()
	case delivery.ErrDeliveryInvalidData:
		status = http.StatusBadRequest
		msg = err.Error()
	case delivery.ErrDeliveryCourierLost:
		status = http.StatusNotFound
		msg = err.Error()
	default:
		status = http.StatusInternalServerError
		msg = delivery.ErrDeliveryDatabase.Error()
	}
	return status, msg
}

func deliveryMapResponse(w http.ResponseWriter, status int, data any, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		stErr, msg := deliveryHandMapError(err)
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
