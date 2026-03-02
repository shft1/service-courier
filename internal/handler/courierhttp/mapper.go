package courierhttp

import (
	"encoding/json"
	"net/http"

	"github.com/shft1/service-courier/internal/domain/courier"
	"github.com/shft1/service-courier/observability/logger"
)

func mapError(err error) (int, string) {
	switch err {
	case courier.ErrCourierInvalidData:
		return http.StatusBadRequest, err.Error()
	case courier.ErrCourierExistPhone:
		return http.StatusConflict, err.Error()
	case courier.ErrCourierEmptyData:
		return http.StatusBadRequest, err.Error()
	case courier.ErrCourierInvalidPhone:
		return http.StatusBadRequest, err.Error()
	case courier.ErrCourierNotFound:
		return http.StatusNotFound, err.Error()
	case courier.ErrCourierInvalidID:
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, courier.ErrCourierDatabase.Error()
	}
}

// writeResponse - подготовка и отправка ответа для клиента
func writeResponse(log logger.Logger, w http.ResponseWriter, status int, data any, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		var msg string
		status, msg = mapError(err)
		if status == http.StatusInternalServerError {
			log.Error("unknown server error", logger.NewField("error", err))
		}
		data = map[string]string{"error": msg}
	}
	w.WriteHeader(status)

	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error("failed to send payload")
	}
}
