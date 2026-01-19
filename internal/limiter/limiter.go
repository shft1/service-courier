package limiter

import (
	"context"
	"sync"
	"time"
)

type tokenBucketLimiter struct {
	mu         sync.Mutex
	capacity   int
	bucket     chan struct{}
	period     time.Duration
	lastRefill time.Time
}

func NewTokenBucketLimiter(period time.Duration, cap int) *tokenBucketLimiter {
	limiter := &tokenBucketLimiter{
		capacity:   cap,
		bucket:     make(chan struct{}, cap),
		period:     period,
		lastRefill: time.Now(),
	}
	for i := 0; i < cap; i++ {
		limiter.bucket <- struct{}{}
	}
	return limiter
}

func (tb *tokenBucketLimiter) StartReplenishment(ctx context.Context) {
	ticker := time.NewTicker(tb.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case tb.bucket <- struct{}{}:
				tb.mu.Lock()
				tb.lastRefill = time.Now()
				tb.mu.Unlock()
			default:
			}
		case <-ctx.Done():
			return
		}
	}
}

func (tb *tokenBucketLimiter) Allow() bool {
	select {
	case <-tb.bucket:
		return true
	default:
		return false
	}
}

func (tb *tokenBucketLimiter) GetRetryAfter() time.Duration {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	nextRefill := tb.lastRefill.Add(tb.period)
	return time.Until(nextRefill)
}

func (tb *tokenBucketLimiter) GetLimit() int {
	return tb.capacity
}
