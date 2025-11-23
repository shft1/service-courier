package router

import (
	courierHandler "service-courier/internal/handler/courier"
	deliveryHandler "service-courier/internal/handler/delivery"
	healthHandler "service-courier/internal/handler/health"
	courierRouter "service-courier/internal/router/courier"
	deliveryRouter "service-courier/internal/router/delivery"
	healthRouter "service-courier/internal/router/health"

	"github.com/go-chi/chi/v5"
)

func SetupRoute(
	hHand *healthHandler.HealthHandler,
	cHand *courierHandler.CourierHandler,
	dHand *deliveryHandler.DeliveryHandler,
) chi.Router {
	mainRouter := chi.NewRouter()
	healthRouter.HealthRoute(mainRouter, hHand)
	courierRouter.CourierRoute(mainRouter, cHand)
	deliveryRouter.DeliveryRoute(mainRouter, dHand)
	return mainRouter
}
