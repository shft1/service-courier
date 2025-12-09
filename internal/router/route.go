package router

import (
	"service-courier/internal/handler/courierhttp"
	"service-courier/internal/handler/deliveryhttp"
	"service-courier/internal/handler/healthhttp"
	"service-courier/internal/router/courierroute"
	"service-courier/internal/router/deliveryroute"
	"service-courier/internal/router/healthroute"

	"github.com/go-chi/chi/v5"
)

// SetupRoute - регистрирует все эндпоинты в роутере
func SetupRoute(
	hHand *healthhttp.HealthHandler,
	cHand *courierhttp.CourierHandler,
	dHand *deliveryhttp.DeliveryHandler,
) chi.Router {
	mainRouter := chi.NewRouter()
	healthroute.HealthRoute(mainRouter, hHand)
	courierroute.CourierRoute(mainRouter, cHand)
	deliveryroute.DeliveryRoute(mainRouter, dHand)
	return mainRouter
}
