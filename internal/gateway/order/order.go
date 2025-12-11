package order

import (
	"context"
	"fmt"
	pb "service-courier/internal/proto/order"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// orderGateway - gRPC шлюз с сервисом доставок
type orderGateway struct {
	client orderClient
}

func NewGateway(clientPB orderClient) *orderGateway {
	return &orderGateway{client: clientPB}
}

// GetOrders - получить заказы через удаленный вызов процедуры
func (og *orderGateway) GetOrders(ctx context.Context, cursor time.Time) ([]*OrderResponse, error) {
	out, err := og.client.GetOrders(ctx, &pb.GetOrdersRequest{From: timestamppb.New(cursor)})
	if err != nil {
		return nil, fmt.Errorf("failed to execute remote procedure: %w", err)
	}
	return toOrderResponseList(out.Orders), nil
}
