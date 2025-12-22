package router

import (
	"net/http"
	"service-courier/internal/handler/courierhttp"
	"service-courier/internal/handler/deliveryhttp"
	"service-courier/internal/handler/healthhttp"
	"service-courier/internal/router/courierroute"
	"service-courier/internal/router/deliveryroute"
	"service-courier/internal/router/healthroute"
	"service-courier/internal/router/metricsroute"

	"github.com/go-chi/chi/v5"
)


type Middleware func(http.Handler) http.Handler


// SetupRoute - регистрирует эндпоинты и middleware в роутере
func SetupRoute(
	metricsMW Middleware,
	hHand *healthhttp.HealthHandler,
	cHand *courierhttp.CourierHandler,
	dHand *deliveryhttp.DeliveryHandler,
	metricsHand http.HandlerFunc,
) chi.Router {
	mainRouter := chi.NewRouter()

	mainRouter.Use(metricsMW)

	healthroute.HealthRoute(mainRouter, hHand)
	courierroute.CourierRoute(mainRouter, cHand)
	deliveryroute.DeliveryRoute(mainRouter, dHand)
	metricsroute.MetricsRoute(mainRouter, metricsHand)

	return mainRouter
}
