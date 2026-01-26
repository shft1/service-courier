package deliveryapp_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"service-courier/internal/service/deliveryapp"
)

func TestFactoryTimeCalculator(t *testing.T) {
	factory := deliveryapp.NewFactoryTimeCalculator()
	unknownErr := fmt.Errorf("some unknown err from factory")

	tests := []struct {
		name          string
		transportType string
		wantType      any
		wantErr       error
	}{
		{
			"valid for foot calc",
			"on_foot",
			deliveryapp.OnFootCalculator{},
			nil,
		},
		{
			"valid for scooter calc",
			"scooter",
			deliveryapp.OnScooterCalculator{},
			nil,
		},
		{
			"valid for car calc",
			"car",
			deliveryapp.OnCarCalculator{},
			nil,
		},
		{
			"unknown transport type",
			"test_transport",
			nil,
			unknownErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := factory.GetDeliveryCalculator(tt.transportType)
			assert.IsType(t, tt.wantType, res)
			if tt.wantErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
