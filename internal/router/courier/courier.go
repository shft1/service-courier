package courier

import (
	"service-courier/internal/handler/courier"

	"github.com/go-chi/chi/v5"
)

func CourierRoute(mr *chi.Mux, hand *courier.CourierHandler) {
	mr.Post("/courier", hand.Create)
	mr.Put("/courier", hand.Update)
	mr.Get("/courier/{id}", hand.GetByID)
	mr.Get("/couriers", hand.GetMulti)
}
