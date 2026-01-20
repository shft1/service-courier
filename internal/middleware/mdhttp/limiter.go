package mdhttp

import (
	"encoding/json"
	"net/http"
	"service-courier/internal/resilience/limiter"
	"strconv"
)

func NewLimiterMiddleware(limiter limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if !limiter.Allow() {
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limiter.GetLimit()))
					w.Header().Set("Retry-After", limiter.GetRetryAfter().String())
					w.WriteHeader(http.StatusTooManyRequests)
					json.NewEncoder(w).Encode(map[string]string{"error": "Rate limit exceeded"})
					return
				}
				next.ServeHTTP(w, r)
			},
		)
	}
}
