package mdhttp

import (
	"encoding/json"
	"net/http"
	"strconv"

	"service-courier/internal/resilience/limiter"
	"service-courier/observability/logger"
)

func NewLimiterMiddleware(log logger.Logger, limiter limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if !limiter.Allow() {
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limiter.GetLimit()))
					w.Header().Set("Retry-After", limiter.GetRetryAfter().String())
					w.WriteHeader(http.StatusTooManyRequests)

					err := json.NewEncoder(w).Encode(map[string]string{"error": "Rate limit exceeded"})
					if err != nil {
						log.Error("limiter middleware: failed to encode response", logger.NewField("error", err))
					}
					return
				}
				next.ServeHTTP(w, r)
			},
		)
	}
}
