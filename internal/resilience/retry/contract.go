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

const IsRetryKey retryContextKey = "is_retry"
