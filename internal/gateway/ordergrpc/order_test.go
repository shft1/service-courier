package ordergrpc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/shft1/service-courier/internal/domain/order"
	"github.com/shft1/service-courier/internal/gateway/ordergrpc"
	"github.com/shft1/service-courier/internal/proto/orderpb"
)

func TestOrderGateway_GetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	createdAt, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00")
	clientErr := errors.New("error from client")
	mockClient := NewMockorderClient(ctrl)
	orderGW := ordergrpc.NewGateway(mockClient)

	tests := []struct {
		name    string
		out     *orderpb.GetOrdersResponse
		err     error
		wantOut []*order.Order
		wantErr error
	}{
		{
			"valid get orders",
			&orderpb.GetOrdersResponse{Orders: []*orderpb.Order{
				{Id: "1", Status: "cooking", CreatedAt: timestamppb.New(createdAt)},
			}},
			nil,
			[]*order.Order{{OrderID: "1", Status: "cooking", CreatedAt: createdAt}},
			nil,
		},
		{
			"wrapped error",
			nil,
			clientErr,
			nil,
			clientErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.EXPECT().
				GetOrders(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(tt.out, tt.err)

			res, err := orderGW.GetOrders(context.Background(), time.Now())

			assert.Equal(t, tt.wantOut, res)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, clientErr)
			}
		})
	}
}

func TestOrderGateway_GetOrderByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	createdAt, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00")
	clientErr := errors.New("error from client")
	mockClient := NewMockorderClient(ctrl)
	orderGW := ordergrpc.NewGateway(mockClient)

	tests := []struct {
		name    string
		out     *orderpb.GetOrderByIdResponse
		err     error
		wantOut *order.Order
		wantErr error
	}{
		{
			"valid get order by ID",
			&orderpb.GetOrderByIdResponse{Order: &orderpb.Order{
				Id: "1", Status: "cooking", CreatedAt: timestamppb.New(createdAt),
			}},
			nil,
			&order.Order{OrderID: "1", Status: "cooking", CreatedAt: createdAt},
			nil,
		},
		{
			"wrapped error",
			nil,
			clientErr,
			nil,
			clientErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.EXPECT().
				GetOrderById(gomock.Any(), gomock.Any()).
				Return(tt.out, tt.err)

			res, err := orderGW.GetOrderByID(context.Background(), order.OrderID{})

			assert.Equal(t, tt.wantOut, res)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, clientErr)
			}
		})
	}
}
