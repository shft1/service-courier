package deliveryhttp

import (
	"encoding/json"
	"net/http"

	"github.com/shft1/service-courier/internal/domain/delivery"
	"github.com/shft1/service-courier/observability/logger"
)

func mapError(err error) (int, string) {
	switch err {
	case delivery.ErrDeliveryExist:
		return http.StatusConflict, err.Error()
	case delivery.ErrDeliveryNotFound:
		return http.StatusNotFound, err.Error()
	case delivery.ErrDeliveryNotAvailableCourier:
		return http.StatusConflict, err.Error()
	case delivery.ErrDeliveryInvalidAssignCourier:
		return http.StatusConflict, err.Error()
	case delivery.ErrDeliveryInvalidOrderID:
		return http.StatusBadRequest, err.Error()
	case delivery.ErrDeliveryEmptyData:
		return http.StatusBadRequest, err.Error()
	case delivery.ErrDeliveryInvalidData:
		return http.StatusBadRequest, err.Error()
	case delivery.ErrDeliveryCourierLost:
		return http.StatusNotFound, err.Error()
	default:
		return http.StatusInternalServerError, delivery.ErrDeliveryDatabase.Error()
	}
}

// writeResponse - подготовка и отправка ответа для клиента
func writeResponse(log logger.Logger, w http.ResponseWriter, status int, data any, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		stErr, msg := mapError(err)
		w.WriteHeader(stErr)

		err = json.NewEncoder(w).Encode(map[string]string{"error": msg})
		if err != nil {
			log.Error("delivery: failed to encode response", logger.NewField("error", err))
		}
		return
	}
	if data == nil {
		w.WriteHeader(status)
		return
	}
	w.WriteHeader(status)

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Error("delivery: failed to encode response", logger.NewField("error", err))
	}
}
