package healthhttp

import (
	"encoding/json"
	"net/http"

	"service-courier/observability/logger"
)

type HealthHandler struct {
	log logger.Logger
}

func NewHealthHandler(log logger.Logger) *HealthHandler {
	return &HealthHandler{log: log}
}

// Ping - проверка ответа от сервиса
func (hh *HealthHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	if err != nil {
		hh.log.Error("health: failed to encode response", logger.NewField("error", err))
	}
}

// HealthCheck - проверка работы сервиса
func (hh *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
