package ordergrpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/shft1/service-courier/internal/proto/orderpb"
)

//go:generate mockgen -source=./contract.go -destination=./mocks_test.go -package=ordergrpc_test

type orderClient interface {
	GetOrders(
		ctx context.Context,
		in *orderpb.GetOrdersRequest,
		opts ...grpc.CallOption,
	) (*orderpb.GetOrdersResponse, error)

	GetOrderById(
		ctx context.Context,
		in *orderpb.GetOrderByIdRequest,
		opts ...grpc.CallOption,
	) (*orderpb.GetOrderByIdResponse, error)
}
