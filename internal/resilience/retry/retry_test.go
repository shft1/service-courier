package retry_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"service-courier/internal/resilience/retry"
)

func TestRetryExecutor_ExecuteWithContext(t *testing.T) {
	exec := retry.NewRetryExecutor(retry.WithMaxAttempts(5))
	fnErr := errors.New("error from test func")
	attempt := 0

	tests := []struct {
		name    string
		fn      func(context.Context) error
		wantErr error
	}{
		{
			"success after retry",
			func(ctx context.Context) error {
				_, ok := ctx.Value(retry.IsRetryKey).(bool)
				if ok && attempt == 3 {
					return nil
				}
				attempt++
				return fnErr
			},
			nil,
		},
		{
			"all attempts fail",
			func(ctx context.Context) error {
				_, ok := ctx.Value(retry.IsRetryKey).(bool)
				if ok && attempt == 6 {
					return nil
				}
				attempt++
				return fnErr
			},
			fnErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := exec.ExecuteWithContext(context.Background(), tt.fn)
			fmt.Println(attempt)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			attempt = 0
		})
	}
}
