package deliveryapp_test

import (
	"context"
	"fmt"
	"service-courier/internal/domain/courier"
	"service-courier/internal/domain/delivery"
	"service-courier/internal/service/deliveryapp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDeliveryService_Assign(t *testing.T) {
	ctrl := gomock.NewController(t)
	mc := NewMockcourierRepository(ctrl)
	md := NewMockdeliveryRepository(ctrl)
	mtx := NewMocktxManagerDo(ctrl)
	mfac := NewMocktimeCalculatorFactory(ctrl)
	mcalc := NewMockTimeCalculator(ctrl)

	tests := []struct {
		name       string
		cRepoAv    *courier.Courier
		cRepoAvErr error
		dRepoCd    *delivery.Delivery
		dRepoCdErr error
		cRepoSbErr error
		srvExpErr  error
		srvExp     *delivery.AssignResult
		input      delivery.OrderID
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
			delivery.OrderID{},
		},
		{
			"delivery exist",
			&courier.Courier{TransportType: "scooter"},
			nil,
			nil,
			delivery.ErrDeliveryExist,
			nil,
			delivery.ErrDeliveryExist,
			nil,
			delivery.OrderID{},
		},
		{
			"unknown error from create_delivery",
			&courier.Courier{TransportType: "scooter"},
			nil,
			nil,
			fmt.Errorf("some unknown wrapped error from repo"),
			nil,
			nil,
			nil,
			delivery.OrderID{},
		},
		{
			"unknown error from set_busy",
			&courier.Courier{TransportType: "scooter"},
			nil,
			&delivery.Delivery{},
			nil,
			fmt.Errorf("some unknown wrapped error from repo"),
			nil,
			nil,
			delivery.OrderID{},
		},
		{
			"courier not found",
			&courier.Courier{TransportType: "scooter"},
			nil,
			&delivery.Delivery{},
			nil,
			courier.ErrCourierNotFound,
			delivery.ErrDeliveryCourierLost,
			nil,
			delivery.OrderID{},
		},
		{
			"valid",
			&courier.Courier{
				ID:            1,
				Name:          "TestName",
				Phone:         "+1234567890",
				Status:        "available",
				TransportType: "scooter",
			},
			nil,
			&delivery.Delivery{
				DeliveryID: 1,
				CourierID:  1,
				OrderID:    "some test orderID",
				AssignedAt: func() time.Time { res, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00"); return res }(),
				Deadline:   func() time.Time { res, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00"); return res }(),
			},
			nil,
			nil,
			nil,
			&delivery.AssignResult{
				CourierID:     1,
				OrderID:       "some test orderID",
				TransportType: "scooter",
				Deadline:      func() time.Time { res, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00"); return res }(),
			},
			delivery.OrderID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc.EXPECT().
				GetAvailable(gomock.Any()).
				Return(tt.cRepoAv, tt.cRepoAvErr)
			if tt.cRepoAvErr == nil {
				mfac.EXPECT().
					GetDeliveryCalculator(gomock.Any()).
					Return(mcalc, nil)
				mcalc.EXPECT().
					Calculate().
					Return(func() time.Time {
						res, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00")
						return res
					}())
				md.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(tt.dRepoCd, tt.dRepoCdErr)
			}
			if tt.cRepoAvErr == nil && tt.dRepoCdErr == nil {
				mc.EXPECT().
					SetBusy(gomock.Any(), gomock.Any()).
					Return(int64(-1), tt.cRepoSbErr)
			}

			mtx.EXPECT().
				Do(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})

			ctx := context.Background()
			s := deliveryapp.NewDeliveryService(md, mc, mfac, mtx)
			res, err := s.Assign(ctx, tt.input)

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
	mc := NewMockcourierRepository(ctrl)
	md := NewMockdeliveryRepository(ctrl)
	mtx := NewMocktxManagerDo(ctrl)
	mfac := NewMocktimeCalculatorFactory(ctrl)

	tests := []struct {
		name        string
		dRepoDel    *delivery.Delivery
		dRepoDelErr error
		cRepoSaErr  error
		srvExpErr   error
		srvExp      *delivery.UnassignResult
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
			&delivery.Delivery{},
			nil,
			fmt.Errorf("some unknown wrapped error from repo"),
			nil,
			nil,
		},
		{
			"courier not found for set_available",
			&delivery.Delivery{},
			nil,
			courier.ErrCourierNotFound,
			delivery.ErrDeliveryCourierLost,
			nil,
		},
		{
			"valid",
			&delivery.Delivery{
				DeliveryID: 1,
				CourierID:  1,
				OrderID:    "some test orderID",
				AssignedAt: func() time.Time { res, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00"); return res }(),
				Deadline:   func() time.Time { res, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00"); return res }(),
			},
			nil,
			nil,
			nil,
			&delivery.UnassignResult{
				OrderID:   "some test orderID",
				Status:    "unassigned",
				CourierID: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md.EXPECT().
				Delete(gomock.Any(), gomock.Any()).
				Return(tt.dRepoDel, tt.dRepoDelErr)
			if tt.dRepoDelErr == nil {
				mc.EXPECT().
					SetAvailable(gomock.Any(), gomock.Any()).
					Return(int64(-1), tt.cRepoSaErr)
			}
			mtx.EXPECT().
				Do(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})
			ctx := context.Background()
			s := deliveryapp.NewDeliveryService(md, mc, mfac, mtx)
			res, err := s.Unassign(ctx, delivery.OrderID{})

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
	mc := NewMockcourierRepository(ctrl)
	md := NewMockdeliveryRepository(ctrl)
	mtx := NewMocktxManagerDo(ctrl)
	mfac := NewMocktimeCalculatorFactory(ctrl)

	tests := []struct {
		name       string
		dRepoRdErr error
	}{
		{"unknown error for recheck_delivery", fmt.Errorf("some unknown wrapped error from repo")},
		{"valid", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc.EXPECT().
				ReleaseStaleBusy(gomock.Any()).
				Return(tt.dRepoRdErr)

			ctx := context.Background()
			s := deliveryapp.NewDeliveryService(md, mc, mfac, mtx)
			err := s.CheckDelivery(ctx)

			assert.ErrorIs(t, err, tt.dRepoRdErr)
		})
	}
}
