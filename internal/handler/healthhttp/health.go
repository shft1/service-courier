package healthhttp

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

func (hh *HealthHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}

func (hh *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
