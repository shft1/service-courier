package mdrpc

import (
	"context"
	"service-courier/observability/logger"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func NewLoggerInterceptor(log logger.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any,
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()

		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(start)

		if err != nil {
			if st, ok := status.FromError(err); !ok {
				log.Warn(
					"unknown error from grpc-request",
					logger.NewField("method", method),
					logger.NewField("error", err),
					logger.NewField("duration", duration.Milliseconds()),
				)
			} else {
				log.Warn(
					"error from grpc-request",
					logger.NewField("method", method),
					logger.NewField("code", st.Code()),
					logger.NewField("msg", st.Message()),
					logger.NewField("duration", duration.Milliseconds()),
				)
			}
		} else {
			log.Info(
				"grpc-request",
				logger.NewField("method", method),
				logger.NewField("duration", duration.Milliseconds()),
			)
		}
		return err
	}
}
