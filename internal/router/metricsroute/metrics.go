package metricsroute

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MetricsRoute - роуты для метрик
func MetricsRoute(mr *chi.Mux, mHand http.HandlerFunc) {
	mr.Get("/metrics", mHand)
}
