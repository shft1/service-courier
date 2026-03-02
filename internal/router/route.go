package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/shft1/service-courier/internal/handler/courierhttp"
	"github.com/shft1/service-courier/internal/handler/deliveryhttp"
	"github.com/shft1/service-courier/internal/handler/healthhttp"
	"github.com/shft1/service-courier/internal/router/courierroute"
	"github.com/shft1/service-courier/internal/router/deliveryroute"
	"github.com/shft1/service-courier/internal/router/healthroute"
	"github.com/shft1/service-courier/internal/router/metricsroute"
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
