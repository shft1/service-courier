package metricsroute

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)


func MetricsRoute(mr *chi.Mux, mHand http.HandlerFunc) {
	mr.Get("/metrics", mHand)
}