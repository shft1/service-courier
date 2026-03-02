package deliveryapp_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/shft1/service-courier/internal/domain/courier"
	"github.com/shft1/service-courier/internal/domain/delivery"
	"github.com/shft1/service-courier/internal/domain/order"
	"github.com/shft1/service-courier/internal/service/deliveryapp"
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
		input      order.OrderID
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
			order.OrderID{},
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
			order.OrderID{},
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
			order.OrderID{},
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
			order.OrderID{},
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
			order.OrderID{},
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
				Deadline:      func() time.Time { t, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00"); return t }(),
			},
			order.OrderID{},
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
			s := deliveryapp.NewDeliveryService(deliveryapp.Arguments{
				DelRepo: md, CourRepo: mc, Factory: mfac, TxManager: mtx,
			})
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
		cRepoSa     int64
		cRepoSaErr  error
		srvExpErr   error
		srvExp      *delivery.UnassignResult
	}{
		{
			"delivery not found",
			nil,
			delivery.ErrDeliveryNotFound,
			int64(-1),
			nil,
			delivery.ErrDeliveryNotFound,
			nil,
		},
		{
			"unknown error for delete_delivery",
			nil,
			fmt.Errorf("some unknown wrapped error from repo"),
			int64(-1),
			nil,
			nil,
			nil,
		},
		{
			"unknown error for set_available",
			&delivery.Delivery{},
			nil,
			int64(-1),
			fmt.Errorf("some unknown wrapped error from repo"),
			nil,
			nil,
		},
		{
			"courier not found for set_available",
			&delivery.Delivery{},
			nil,
			int64(-1),
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
			int64(1),
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
			mtx.EXPECT().
				Do(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})

			md.EXPECT().
				Delete(gomock.Any(), gomock.Any()).
				Return(tt.dRepoDel, tt.dRepoDelErr)

			if tt.dRepoDelErr == nil {
				mc.EXPECT().
					SetAvailable(gomock.Any(), gomock.Any()).
					Return(tt.cRepoSa, tt.cRepoSaErr)
			}
			ctx := context.Background()
			s := deliveryapp.NewDeliveryService(deliveryapp.Arguments{
				DelRepo: md, CourRepo: mc, Factory: mfac, TxManager: mtx,
			})
			res, err := s.Unassign(ctx, order.OrderID{})

			assert.Equal(t, tt.srvExp, res)
			if err != nil && tt.srvExpErr == nil {
				assert.Contains(t, err.Error(), "service: failed to work with delivery")
			} else {
				assert.ErrorIs(t, err, tt.srvExpErr)
			}
		})
	}
}

func TestDeliveryService_Complete(t *testing.T) {
	ctrl := gomock.NewController(t)
	md := NewMockdeliveryRepository(ctrl)
	mc := NewMockcourierRepository(ctrl)
	mfac := NewMocktimeCalculatorFactory(ctrl)
	mtx := NewMocktxManagerDo(ctrl)
	unknownErr := fmt.Errorf("some unknown wrapped error from repo")

	delSrv := deliveryapp.NewDeliveryService(deliveryapp.Arguments{
		DelRepo: md, CourRepo: mc, Factory: mfac, TxManager: mtx,
	})

	tests := []struct {
		name        string
		dRepoGet    *delivery.Delivery
		dRepoGetErr error
		cRepoSa     int64
		cRepoSaErr  error
		srvExp      *delivery.CompleteResult
		srvExpErr   error
	}{
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
			int64(1),
			nil,
			&delivery.CompleteResult{
				CourierID: int64(1),
				OrderID:   "some test orderID",
				Deadline:  func() time.Time { res, _ := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00"); return res }(),
			},
			nil,
		},
		{
			"delivery not found",
			nil,
			delivery.ErrDeliveryNotFound,
			int64(-1),
			nil,
			nil,
			delivery.ErrDeliveryNotFound,
		},
		{
			"courier not found",
			nil,
			courier.ErrCourierNotFound,
			int64(-1),
			nil,
			nil,
			delivery.ErrDeliveryCourierLost,
		},
		{
			"wrapped unknown error from delivery.Get",
			nil,
			unknownErr,
			int64(-1),
			nil,
			nil,
			unknownErr,
		},
		{
			"wrapped unknown error from courier.SetAvailable",
			&delivery.Delivery{},
			nil,
			int64(-1),
			unknownErr,
			nil,
			unknownErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mtx.EXPECT().
				Do(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})

			md.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				Return(tt.dRepoGet, tt.dRepoGetErr)

			if tt.dRepoGetErr == nil {
				mc.EXPECT().
					SetAvailable(gomock.Any(), gomock.Any()).
					Return(tt.cRepoSa, tt.cRepoSaErr)
			}
			res, err := delSrv.Complete(context.Background(), order.OrderID{})

			assert.Equal(t, tt.srvExp, res)
			assert.ErrorIs(t, err, tt.srvExpErr)
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
			s := deliveryapp.NewDeliveryService(deliveryapp.Arguments{
				DelRepo: md, CourRepo: mc, Factory: mfac, TxManager: mtx,
			})
			err := s.CheckDelivery(ctx)

			assert.ErrorIs(t, err, tt.dRepoRdErr)
		})
	}
}
