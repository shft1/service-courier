package router

import (
	"service-courier/internal/handler"

	"github.com/go-chi/chi/v5"
)

func HealthRoute(mr *chi.Mux, hh *handler.HealthHandler) {
	mr.Get("/ping", hh.Ping)
	mr.Head("/healthcheck", hh.HealthCheck)
}
