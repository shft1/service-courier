package delivery

import (
	"context"
	"fmt"
	"service-courier/internal/entity/courier"
	"service-courier/internal/entity/delivery"
	"service-courier/internal/service/delivery/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDeliveryService_Assign(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mocks.NewMockcourierRepository(ctrl)
	md := mocks.NewMockdeliveryRepository(ctrl)
	mtx := mocks.NewMockTxManagerDo(ctrl)

	tests := []struct {
		name       string
		cRepoAv    *courier.CourierGet
		cRepoAvErr error
		dRepoCd    *delivery.DeliveryGet
		dRepoCdErr error
		cRepoSbErr error
		srvExpErr  error
		srvExp     *delivery.DeliveryAssign
		input      *delivery.DeliveryOrderID
	}{
		{
			"busy couriers",
			nil,
			courier.ErrCourierAvailable,
			nil,
			nil,
			nil,
			delivery.ErrDeliveryNotAvailableCourier,
			nil,
			&delivery.DeliveryOrderID{},
		},
		{
			"delivery exist",
			&courier.CourierGet{TransportType: "scooter"},
			nil,
			nil,
			delivery.ErrDeliveryExist,
			nil,
			delivery.ErrDeliveryExist,
			nil,
			&delivery.DeliveryOrderID{},
		},
		{
			"unknown error from create_delivery",
			&courier.CourierGet{TransportType: "scooter"},
			nil,
			nil,
			fmt.Errorf("some unknown wrapped error from repo"),
			nil,
			nil,
			nil,
			&delivery.DeliveryOrderID{},
		},
		{
			"unknown error from set_busy",
			&courier.CourierGet{TransportType: "scooter"},
			nil,
			&delivery.DeliveryGet{},
			nil,
			fmt.Errorf("some unknown wrapped error from repo"),
			nil,
			nil,
			&delivery.DeliveryOrderID{},
		},
		{
			"courier not found",
			&courier.CourierGet{TransportType: "scooter"},
			nil,
			&delivery.DeliveryGet{},
			nil,
			courier.ErrCourierNotFound,
			delivery.ErrDeliveryCourierLost,
			nil,
			&delivery.DeliveryOrderID{},
		},
		{
			"valid",
			&courier.CourierGet{
				ID:            1,
				Name:          "TestName",
				Phone:         "+1234567890",
				Status:        "available",
				TransportType: "scooter",
			},
			nil,
			&delivery.DeliveryGet{
				DeliveryID:       1,
				CourierID:        1,
				OrderID:          "some test orderID",
				AssignedAt:       func() time.Time { res, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00"); return res }(),
				DeliveryDeadline: func() time.Time { res, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00"); return res }(),
			},
			nil,
			nil,
			nil,
			&delivery.DeliveryAssign{
				CourierID:        1,
				OrderID:          "some test orderID",
				TransportType:    "scooter",
				DeliveryDeadline: func() time.Time { res, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00"); return res }(),
			},
			&delivery.DeliveryOrderID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc.EXPECT().
				GetAvailable(gomock.Any()).
				Return(tt.cRepoAv, tt.cRepoAvErr)
			if tt.cRepoAvErr == nil {
				md.EXPECT().
					CreateDelivery(gomock.Any(), gomock.Any()).
					Return(tt.dRepoCd, tt.dRepoCdErr)
			}
			if tt.cRepoAvErr == nil && tt.dRepoCdErr == nil {
				mc.EXPECT().
					SetBusy(gomock.Any(), gomock.Any()).
					Return(tt.cRepoSbErr)
			}
			mtx.EXPECT().
				Do(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})
			ctx := context.Background()
			s := NewDeliveryService(md, mc, mtx)
			res, err := s.DeliveryAssign(ctx, tt.input)

			assert.Equal(t, tt.srvExp, res)
			if err != nil && tt.srvExpErr == nil {
				assert.Contains(t, err.Error(), "service: failed to work with delivery")
			} else {
				assert.ErrorIs(t, err, tt.srvExpErr)
			}
		})
	}
}

func TestDeliveryService_Unassign(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mocks.NewMockcourierRepository(ctrl)
	md := mocks.NewMockdeliveryRepository(ctrl)
	mtx := mocks.NewMockTxManagerDo(ctrl)

	tests := []struct {
		name        string
		dRepoDel    *delivery.DeliveryGet
		dRepoDelErr error
		cRepoSaErr  error
		srvExpErr   error
		srvExp      *delivery.DeliveryUnassign
	}{
		{
			"delivery not found",
			nil,
			delivery.ErrDeliveryNotFound,
			nil,
			delivery.ErrDeliveryNotFound,
			nil,
		},
		{
			"unknown error for delete_delivery",
			nil,
			fmt.Errorf("some unknown wrapped error from repo"),
			nil,
			nil,
			nil,
		},
		{
			"unknown error for set_available",
			&delivery.DeliveryGet{},
			nil,
			fmt.Errorf("some unknown wrapped error from repo"),
			nil,
			nil,
		},
		{
			"courier not found for set_available",
			&delivery.DeliveryGet{},
			nil,
			courier.ErrCourierNotFound,
			delivery.ErrDeliveryCourierLost,
			nil,
		},
		{
			"valid",
			&delivery.DeliveryGet{
				DeliveryID:       1,
				CourierID:        1,
				OrderID:          "some test orderID",
				AssignedAt:       func() time.Time { res, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00"); return res }(),
				DeliveryDeadline: func() time.Time { res, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00"); return res }(),
			},
			nil,
			nil,
			nil,
			&delivery.DeliveryUnassign{
				OrderID:   "some test orderID",
				Status:    "unassigned",
				CourierID: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md.EXPECT().
				DeleteDelivery(gomock.Any(), gomock.Any()).
				Return(tt.dRepoDel, tt.dRepoDelErr)
			if tt.dRepoDelErr == nil {
				mc.EXPECT().
					SetAvailable(gomock.Any(), gomock.Any()).
					Return(tt.cRepoSaErr)
			}
			mtx.EXPECT().
				Do(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})
			ctx := context.Background()
			s := NewDeliveryService(md, mc, mtx)
			res, err := s.DeliveryUnassign(ctx, nil)

			assert.Equal(t, tt.srvExp, res)
			if err != nil && tt.srvExpErr == nil {
				assert.Contains(t, err.Error(), "service: failed to work with delivery")
			} else {
				assert.ErrorIs(t, err, tt.srvExpErr)
			}
		})
	}
}

func TestDeliveryService_Check(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := mocks.NewMockcourierRepository(ctrl)
	md := mocks.NewMockdeliveryRepository(ctrl)
	mtx := mocks.NewMockTxManagerDo(ctrl)

	tests := []struct {
		name       string
		dRepoRdErr error
	}{
		{"unknown error for recheck_delivery", fmt.Errorf("some unknown wrapped error from repo")},
		{"valid", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md.EXPECT().
				RecheckDelivery(gomock.Any()).
				Return(tt.dRepoRdErr)

			ctx := context.Background()
			s := NewDeliveryService(md, mc, mtx)
			err := s.DeliveryCheck(ctx)

			assert.ErrorIs(t, err, tt.dRepoRdErr)
		})
	}
}
