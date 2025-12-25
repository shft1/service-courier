package ordergrpc

import (
	"context"

	"service-courier/internal/proto/orderpb"

	"google.golang.org/grpc"
)

type orderClient interface {
	GetOrders(ctx context.Context, in *orderpb.GetOrdersRequest, opts ...grpc.CallOption) (*orderpb.GetOrdersResponse, error)
	GetOrderById(ctx context.Context, in *orderpb.GetOrderByIdRequest, opts ...grpc.CallOption) (*orderpb.GetOrderByIdResponse, error)
}
