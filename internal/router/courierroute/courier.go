package courierroute

import (
	"service-courier/internal/handler/courierhttp"

	"github.com/go-chi/chi/v5"
)

// CourierRoute - роуты для курьеров
func CourierRoute(mr *chi.Mux, handler *courierhttp.CourierHandler) {
	mr.Post("/courier", handler.Create)
	mr.Put("/courier", handler.Update)
	mr.Get("/courier/{id}", handler.GetByID)
	mr.Get("/couriers", handler.GetMulti)
}
