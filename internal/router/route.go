package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoute(pool *pgxpool.Pool) chi.Router {
	mainRouter := chi.NewRouter()
	HealthRoute(mainRouter)
	CourierRoute(mainRouter, pool)
	return mainRouter
}
