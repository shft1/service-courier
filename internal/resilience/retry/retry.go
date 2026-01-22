package retry

import (
	"context"
	"fmt"
	"time"

	"service-courier/internal/domain/order"
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
		retryExecutor.strategy = NewExponentialBackoffWithJitter(Arguments{
			Multi: 2.0, Jitter: 0.1, InitDelay: 100 * time.Millisecond, MaxDelay: 5 * time.Second,
		})
	}
	if retryExecutor.shouldRetry == nil {
		retryExecutor.shouldRetry = func(err error) bool { return err != nil }
	}
	return retryExecutor
}

func (r *retryExecutor) ExecuteWithContext(ctx context.Context, fn func(context.Context) error) error {
	var lastErr error

	for attempt := 0; attempt < r.maxAttempts; attempt++ {
		isRetry := attempt > 0
		attemptCtx := context.WithValue(ctx, isRetryKey, isRetry)

		err := fn(attemptCtx)
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

func (r *retryExecutor) IsRetryFromContext(ctx context.Context) bool {
	if isRetry, ok := ctx.Value(isRetryKey).(bool); ok {
		return isRetry
	}
	return false
}
