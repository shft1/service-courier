package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTPMetrics - HTTP-метрики приложения
type HTTPMetrics struct {
	Request *prometheus.CounterVec
	Duration *prometheus.HistogramVec
}

// NewHTTPMetrics - создание и регистрация HTTP-метрик приложения
func NewHTTPMetrics() *HTTPMetrics {
	return &HTTPMetrics{
		Request: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},[]string{"method", "path", "status"}),
		Duration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests in seconds",
			Buckets: []float64{0.005, 0.1, 0.5, 0.75, 1, 5},
		}, []string{"method", "path", "status"}),
	}
}
