package mdhttp

import (
	"net/http"
	"service-courier/observability/metrics/metricshttp"
	"strconv"
	"time"
)

// NewMetricsMiddleware - конструктор Middleware для метрик
func NewMetricsMiddleware(m *metricshttp.HTTPMetrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()

				rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

				next.ServeHTTP(rw, r)

				duration := time.Since(start).Seconds()
				statusCode := strconv.Itoa(rw.statusCode)

				m.Duration.WithLabelValues(r.Method, r.URL.Path, statusCode).Observe(duration)
				m.Request.WithLabelValues(r.Method, r.URL.Path, statusCode).Inc()

				if statusCode == strconv.Itoa(http.StatusTooManyRequests) {
					m.RateLimit.WithLabelValues(statusCode).Inc()
				}
			},
		)
	}
}
