package mdrpc

import (
	"context"
	"service-courier/internal/resilience/retry"

	"google.golang.org/grpc"
)

func NewRetryInterceptor(retry retry.Retry) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any,
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return retry.ExecuteWithContext(ctx, func(ctx context.Context) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		})
	}
}
