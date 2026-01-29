package limiter

import "time"

type RateLimiter interface {
	Allow() bool
	GetLimit() int
	GetRetryAfter() time.Duration
}
