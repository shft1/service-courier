package courier

import (
	"encoding/json"
	"net/http"
	"service-courier/internal/entity/courier"
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
		msg = courier.ErrCourierDatabase.Error()
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
