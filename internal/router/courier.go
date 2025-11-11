package router

import (
	"service-courier/internal/handler"

	"github.com/go-chi/chi/v5"
)

func CourierRoute(mr *chi.Mux, hand *handler.CourierHandler) {
	mr.Post("/courier", hand.Create)
	mr.Put("/courier", hand.Update)
	mr.Get("/courier/{id}", hand.GetByID)
	mr.Get("/couriers", hand.GetMulti)
}
