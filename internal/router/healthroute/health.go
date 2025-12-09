package healthroute

import (
	"service-courier/internal/handler/healthhttp"

	"github.com/go-chi/chi/v5"
)

// HealthRoute - роуты для health-проверок
func HealthRoute(mr *chi.Mux, handler *healthhttp.HealthHandler) {
	mr.Get("/ping", handler.Ping)
	mr.Head("/healthcheck", handler.HealthCheck)
}
