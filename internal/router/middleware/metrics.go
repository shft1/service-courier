package middleware

import (
	"net/http"
	"service-courier/observability/metrics"
	"strconv"
	"time"
)


func NewMetricsMiddleware(m *metrics.HTTPMetrics) func(http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(
			func (w http.ResponseWriter, r *http.Request) {
				start := time.Now()

				rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

				next.ServeHTTP(rw, r)

				duration := time.Since(start).Seconds()
				statusCode := strconv.Itoa(rw.statusCode)

				m.Duration.WithLabelValues(r.Method, r.URL.Path, statusCode).Observe(duration)
				m.Request.WithLabelValues(r.Method, r.URL.Path, statusCode).Inc()
			},
		)
	}
}