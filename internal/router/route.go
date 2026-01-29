package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"service-courier/internal/handler/courierhttp"
	"service-courier/internal/handler/deliveryhttp"
	"service-courier/internal/handler/healthhttp"
	"service-courier/internal/router/courierroute"
	"service-courier/internal/router/deliveryroute"
	"service-courier/internal/router/healthroute"
	"service-courier/internal/router/metricsroute"
)

type Middleware func(http.Handler) http.Handler

// SetupRoute - регистрирует эндпоинты и middleware в роутере
func SetupRoute(
	loggerMW Middleware,
	metricsMW Middleware,
	limiter Middleware,
	hHand *healthhttp.HealthHandler,
	cHand *courierhttp.CourierHandler,
	dHand *deliveryhttp.DeliveryHandler,
	metricsHand http.HandlerFunc,
) chi.Router {
	mainRouter := chi.NewRouter()

	mainRouter.Use(loggerMW)
	mainRouter.Use(metricsMW)
	mainRouter.Use(limiter)

	healthroute.HealthRoute(mainRouter, hHand)
	courierroute.CourierRoute(mainRouter, cHand)
	deliveryroute.DeliveryRoute(mainRouter, dHand)
	metricsroute.MetricsRoute(mainRouter, metricsHand)

	return mainRouter
}
