package retry

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ShouldRetry(err error) bool {
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.Unavailable:
			return true
		case codes.DeadlineExceeded:
			return true
		case codes.ResourceExhausted:
			return true
		}
	}
	return false
}
