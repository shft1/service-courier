package middleware

import (
	"net/http"
	"service-courier/observability/logger"
	"time"
)

// NewLoggerMiddleware - конструктор Middleware для логгирования
func NewLoggerMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()

				rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

				next.ServeHTTP(rw, r)

				duration := time.Since(start)
				statusCode := rw.statusCode

				log.Info(
					"http-request",
					logger.NewField("method", r.Method),
					logger.NewField("path", r.URL.Path),
					logger.NewField("status", statusCode),
					logger.NewField("duration", duration.Milliseconds()),
				)
				if statusCode == http.StatusTooManyRequests {
					log.Info(
						"rate limit exceeded",
						logger.NewField("method", r.Method),
						logger.NewField("path", r.URL.Path),
						logger.NewField("status", statusCode),
						logger.NewField("duration", duration.Milliseconds()),
					)
				}
			},
		)
	}
}
