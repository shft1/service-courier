package retry

import (
	"math"
	"math/rand"
	"time"
)

type Arguments struct {
	Multi     float64
	Jitter    float64
	InitDelay time.Duration
	MaxDelay  time.Duration
}

type exponentialBackoffWithJitter struct {
	multiplier   float64
	jitterFactor float64
	initialDelay time.Duration
	maxDelay     time.Duration
}

func NewExponentialBackoffWithJitter(args Arguments) *exponentialBackoffWithJitter {
	return &exponentialBackoffWithJitter{
		multiplier:   args.Multi,
		jitterFactor: args.Jitter,
		initialDelay: args.InitDelay,
		maxDelay:     args.MaxDelay,
	}
}

func (e *exponentialBackoffWithJitter) NextDelay(attempts int) time.Duration {
	delay := float64(e.initialDelay) * math.Pow(e.multiplier, float64(attempts))
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
