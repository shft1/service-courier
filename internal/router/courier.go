package router

import (
	"service-courier/internal/handler"
	"service-courier/internal/repository"
	"service-courier/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CourierRoute(mr *chi.Mux, pool *pgxpool.Pool) {
	cr := repository.NewCourierRepository(pool)
	cs := service.NewCourierService(cr)
	ch := handler.NewCourierHandler(cs)
	mr.Post("/courier", ch.Create)
	mr.Put("/courier", ch.Update)
	mr.Get("/courier/{id}", ch.GetByID)
	mr.Get("/couriers", ch.GetMulti)
}
