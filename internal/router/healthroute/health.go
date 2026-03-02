package healthroute

import (
	"github.com/go-chi/chi/v5"

	"github.com/shft1/service-courier/internal/handler/healthhttp"
)

// HealthRoute - роуты для health-проверок
func HealthRoute(mr *chi.Mux, handler *healthhttp.HealthHandler) {
	mr.Get("/ping", handler.Ping)
	mr.Head("/healthcheck", handler.HealthCheck)
}
