package health

import (
	"service-courier/internal/handler/health"

	"github.com/go-chi/chi/v5"
)

func HealthRoute(mr *chi.Mux, hh *health.HealthHandler) {
	mr.Get("/ping", hh.Ping)
	mr.Head("/healthcheck", hh.HealthCheck)
}
