package router

import (
	courierHandler "service-courier/internal/handler/courier"
	healthHandler "service-courier/internal/handler/health"
	courierRoute "service-courier/internal/router/courier"
	healthRouter "service-courier/internal/router/health"

	"github.com/go-chi/chi/v5"
)

func SetupRoute(hHand *healthHandler.HealthHandler, crHand *courierHandler.CourierHandler) chi.Router {
	mainRouter := chi.NewRouter()
	healthRouter.HealthRoute(mainRouter, hHand)
	courierRoute.CourierRoute(mainRouter, crHand)
	return mainRouter
}
