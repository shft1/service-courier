package retry

import (
	"context"
	"fmt"
	"service-courier/internal/domain/order"
	"time"
)

type retryExecutor struct {
	maxAttempts int
	strategy    strategy
	shouldRetry func(error) bool
}

func NewRetryExecutor(opts ...option) *retryExecutor {
	retryExecutor := &retryExecutor{}
	for _, opt := range opts {
		opt(retryExecutor)
	}
	if retryExecutor.maxAttempts == 0 {
		retryExecutor.maxAttempts = 3
	}
	if retryExecutor.strategy == nil {
		retryExecutor.strategy = NewExponentialBackoffWithJitter(2.0, 0.1, 100*time.Millisecond, 5*time.Second)
	}
	if retryExecutor.shouldRetry == nil {
		retryExecutor.shouldRetry = func(err error) bool { return err != nil }
	}
	return retryExecutor
}

func (r *retryExecutor) ExecuteWithContext(ctx context.Context, fn func(context.Context) error) error {
	var lastErr error
	for attempt := 1; attempt <= r.maxAttempts; attempt++ {
		err := fn(ctx)
		if err == nil {
			return nil
		}
		if !r.shouldRetry(err) {
			return err
		}
		delay := r.strategy.NextDelay(attempt)
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return fmt.Errorf("%w: %v", order.ErrServiceUnavailable, lastErr)
}
