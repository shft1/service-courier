package limiter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"service-courier/internal/resilience/limiter"
)

func TestTokenBucketLimiter_Allow(t *testing.T) {
	tests := []struct {
		name    string
		cap     int
		wantRes bool
	}{
		{
			"allow with tokens",
			1,
			true,
		},
		{
			"reject without tokens",
			0,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bucketLim := limiter.NewTokenBucketLimiter(time.Millisecond, tt.cap)
			isAllow := bucketLim.Allow()
			assert.Equal(t, tt.wantRes, isAllow)
		})
	}
}
