package retry

import (
	"context"
	"time"
)

type Retry interface {
	ExecuteWithContext(context.Context, func(context.Context) error) error
}

type strategy interface {
	NextDelay(attempt int) time.Duration
}

type option func(*retryExecutor)

type retryContextKey string

const isRetryKey retryContextKey = "is_retry"
