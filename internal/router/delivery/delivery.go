package delivery

import (
	"service-courier/internal/handler/delivery"

	"github.com/go-chi/chi/v5"
)

func DeliveryRoute(mr *chi.Mux, handler *delivery.DeliveryHandler) {
	mr.Post("/delivery/assign", handler.DeliveryAssign)
	mr.Post("/delivery/unassign", handler.DeliveryUnassign)
}
