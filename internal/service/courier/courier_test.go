package courier_test

import (
	"context"
	"fmt"
	"service-courier/internal/entity/courier"
	courierservice "service-courier/internal/service/courier"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCourierService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockcourierRepository(ctrl)

	tests := []struct {
		name    string
		repoErr error
	}{
		{"phone exist", courier.ErrCourierExistPhone},
		{"unknown error", fmt.Errorf("some unknown wrapped error from repo")},
		{"valid", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(tt.repoErr)

			ctx := context.Background()
			courier := &courier.CourierCreate{}
			s := courierservice.NewCourierService(m)
			err := s.Create(ctx, courier)

			assert.ErrorIs(t, err, tt.repoErr)
		})
	}
}

func TestCourierService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockcourierRepository(ctrl)

	tests := []struct {
		name    string
		repoErr error
	}{
		{"phone exist", courier.ErrCourierExistPhone},
		{"courier not found", courier.ErrCourierNotFound},
		{"unknown error", fmt.Errorf("some unknown wrapped error from repo")},
		{"valid", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.EXPECT().
				Update(gomock.Any(), gomock.Any()).
				Return(tt.repoErr)

			ctx := context.Background()
			courier := &courier.CourierUpdate{}
			s := courierservice.NewCourierService(m)
			err := s.Update(ctx, courier)

			assert.ErrorIs(t, err, tt.repoErr)
		})
	}
}

func TestCourierService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockcourierRepository(ctrl)

	tests := []struct {
		name    string
		repoRes *courier.CourierGet
		repoErr error
		srvRes  *courier.CourierGet
	}{
		{"courier not found", nil, courier.ErrCourierNotFound, nil},
		{"unknown error", nil, fmt.Errorf("some unknown wrapped error from repo"), nil},
		{"valid", &courier.CourierGet{}, nil, &courier.CourierGet{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.EXPECT().
				GetByID(gomock.Any(), gomock.Any()).
				Return(tt.repoRes, tt.repoErr)

			ctx, id := context.Background(), 1
			s := courierservice.NewCourierService(m)
			res, err := s.GetByID(ctx, id)

			assert.Equal(t, res, tt.srvRes)
			assert.ErrorIs(t, err, tt.repoErr)
		})
	}
}

func TestCourierService_GetMulti(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockcourierRepository(ctrl)

	tests := []struct {
		name    string
		repoRes []courier.CourierGet
		repoErr error
		srvRes  []courier.CourierGet
	}{
		{"unknown error", nil, fmt.Errorf("some unknown wrapped error from repo"), nil},
		{"valid", []courier.CourierGet{{}}, nil, []courier.CourierGet{{}}},
	}

	for _, tt := range tests {
		m.EXPECT().
			GetMulti(gomock.Any()).
			Return(tt.repoRes, tt.repoErr)

		ctx := context.Background()
		s := courierservice.NewCourierService(m)
		res, err := s.GetMulti(ctx)

		assert.Equal(t, res, tt.srvRes)
		assert.ErrorIs(t, err, tt.repoErr)
	}
}
