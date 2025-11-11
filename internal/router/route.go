package router

import (
	"service-courier/internal/handler"

	"github.com/go-chi/chi/v5"
)

func SetupRoute(hHand *handler.HealthHandler, crHand *handler.CourierHandler) chi.Router {
	mainRouter := chi.NewRouter()
	HealthRoute(mainRouter, hHand)
	CourierRoute(mainRouter, crHand)
	return mainRouter
}
