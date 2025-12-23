package healthhttp

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}


func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Ping - проверка ответа от сервиса
func (hh *HealthHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}

// HealthCheck - проверка работы сервиса
func (hh *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
