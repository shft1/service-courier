package retry

import "time"

func WithMaxAttempts(attempts int) option {
	return func(r *retryExecutor) {
		if attempts <= 0 {
			attempts = 3
		}
		r.maxAttempts = attempts
	}
}

func WithStrategy(strategy strategy) option {
	return func(r *retryExecutor) {
		if strategy == nil {
			strategy = NewExponentialBackoffWithJitter(2.0, 0.1, 100*time.Millisecond, 5*time.Second)
		}
		r.strategy = strategy
	}
}

func WithShouldRetry(fn func(error) bool) option {
	return func(r *retryExecutor) {
		if fn == nil {
			fn = func(err error) bool { return err != nil }
		}
		r.shouldRetry = fn
	}
}
