package courierapp_test

import (
	"context"
	"fmt"
	"service-courier/internal/domain/courier"
	"service-courier/internal/service/courierapp"
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
				Return(int64(-1), tt.repoErr)

			ctx := context.Background()
			courier := &courier.CourierCreate{}
			s := courierapp.NewCourierService(m)
			_, err := s.Create(ctx, courier)

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
				Return(int64(-1), tt.repoErr)

			ctx := context.Background()
			courier := &courier.CourierUpdate{}
			s := courierapp.NewCourierService(m)
			_, err := s.Update(ctx, courier)

			assert.ErrorIs(t, err, tt.repoErr)
		})
	}
}

func TestCourierService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockcourierRepository(ctrl)

	tests := []struct {
		name    string
		repoRes *courier.Courier
		repoErr error
		srvRes  *courier.Courier
	}{
		{"courier not found", nil, courier.ErrCourierNotFound, nil},
		{"unknown error", nil, fmt.Errorf("some unknown wrapped error from repo"), nil},
		{"valid", &courier.Courier{}, nil, &courier.Courier{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.EXPECT().
				GetByID(gomock.Any(), gomock.Any()).
				Return(tt.repoRes, tt.repoErr)

			ctx, id := context.Background(), 1
			s := courierapp.NewCourierService(m)
			res, err := s.GetByID(ctx, int64(id))

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
		repoRes []courier.Courier
		repoErr error
		srvRes  []courier.Courier
	}{
		{"unknown error", nil, fmt.Errorf("some unknown wrapped error from repo"), nil},
		{"valid", []courier.Courier{{}}, nil, []courier.Courier{{}}},
	}

	for _, tt := range tests {
		m.EXPECT().
			GetMulti(gomock.Any()).
			Return(tt.repoRes, tt.repoErr)

		ctx := context.Background()
		s := courierapp.NewCourierService(m)
		res, err := s.GetMulti(ctx)

		assert.Equal(t, res, tt.srvRes)
		assert.ErrorIs(t, err, tt.repoErr)
	}
}
