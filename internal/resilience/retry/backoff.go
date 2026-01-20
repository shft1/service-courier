package retry

import (
	"math"
	"math/rand/v2"
	"time"
)

type exponentialBackoffWithJitter struct {
	multiplier   float64
	jitterFactor float64
	initialDelay time.Duration
	maxDelay     time.Duration
}

func NewExponentialBackoffWithJitter(multi, jitter float64, initDelay, maxDelay time.Duration) *exponentialBackoffWithJitter {
	return &exponentialBackoffWithJitter{
		multiplier:   multi,
		jitterFactor: jitter,
		initialDelay: initDelay,
		maxDelay:     maxDelay,
	}
}

func (e *exponentialBackoffWithJitter) NextDelay(attempts int) time.Duration {
	delay := float64(e.initialDelay) * math.Pow(e.multiplier, float64(attempts-1))
	if delay > float64(e.maxDelay) {
		delay = float64(e.maxDelay)
	}
	jitter := delay * e.jitterFactor * (rand.Float64()*2 - 1)
	delay += jitter
	if delay < 0 {
		delay = 0
	}
	return time.Duration(delay)
}
