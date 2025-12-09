package courierhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"service-courier/internal/domain/courier"
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

func writeResponse(w http.ResponseWriter, status int, data any, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		var msg string
		status, msg = mapError(err)
		data = map[string]string{"error": msg}
	}
	w.WriteHeader(status)

	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Println("failed to send payload")
	}
}
