package ordergrpc

import (
	"context"
	"fmt"
	"service-courier/internal/domain/order"
	"service-courier/internal/proto/orderpb"
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
func (og *orderGateway) GetOrders(ctx context.Context, cursor time.Time) ([]*order.Order, error) {
	out, err := og.client.GetOrders(ctx, &orderpb.GetOrdersRequest{From: timestamppb.New(cursor)})
	if err != nil {
		return nil, fmt.Errorf("failed to execute remote get orders: %w", err)
	}
	return toDomainOrderList(out.Orders), nil
}

// GetOrderByID - получить заказ по ID через удаленный вызов процедуры
func (og *orderGateway) GetOrderByID(ctx context.Context, orderID order.OrderID) (*order.Order, error) {
	out, err := og.client.GetOrderById(ctx, &orderpb.GetOrderByIdRequest{Id: orderID.OrderID})
	if err != nil {
		return nil, fmt.Errorf("failed to execute remote get order: %w", err)
	}
	return toDomainOrder(out.Order), nil
}
