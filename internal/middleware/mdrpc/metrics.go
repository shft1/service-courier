package mdrpc

import (
	"context"
	"service-courier/observability/metrics/metricsrpc"

	"google.golang.org/grpc"
)

func NewMetricsInterceptor(m *metricsrpc.RPCMetrics, isRetry func(context.Context) bool) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any,
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if isRetry(ctx) {
			m.Retry.WithLabelValues("retry").Inc()
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
