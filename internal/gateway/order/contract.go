package order

import (
	"context"

	pb "service-courier/internal/proto/order"

	"google.golang.org/grpc"
)

type orderClient interface {
	GetOrders(ctx context.Context, in *pb.GetOrdersRequest, opts ...grpc.CallOption) (*pb.GetOrdersResponse, error)
}
