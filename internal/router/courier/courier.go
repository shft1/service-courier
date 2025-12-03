package courier

import (
	"service-courier/internal/handler/courier"

	"github.com/go-chi/chi/v5"
)

func CourierRoute(mr *chi.Mux, handler *courier.CourierHandler) {
	mr.Post("/courier", handler.Create)
	mr.Put("/courier", handler.Update)
	mr.Get("/courier/{id}", handler.GetByID)
	mr.Get("/couriers", handler.GetMulti)
}
