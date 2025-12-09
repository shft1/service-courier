package deliveryroute

import (
	"service-courier/internal/handler/deliveryhttp"

	"github.com/go-chi/chi/v5"
)

// DeliveryRoute - роуты для доставок
func DeliveryRoute(mr *chi.Mux, handler *deliveryhttp.DeliveryHandler) {
	mr.Post("/delivery/assign", handler.Assign)
	mr.Post("/delivery/unassign", handler.Unassign)
}
