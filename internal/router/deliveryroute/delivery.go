package deliveryroute

import (
	"github.com/go-chi/chi/v5"

	"service-courier/internal/handler/deliveryhttp"
)

// DeliveryRoute - роуты для доставок
func DeliveryRoute(mr *chi.Mux, handler *deliveryhttp.DeliveryHandler) {
	mr.Post("/delivery/assign", handler.Assign)
	mr.Post("/delivery/unassign", handler.Unassign)
}
