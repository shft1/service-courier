package deliveryapp_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/shft1/service-courier/internal/service/deliveryapp"
)

func TestFactoryEventStrategy_GetEventStrategy(t *testing.T) {
	ctrl := gomock.NewController(t)
	executor := NewMockdeliveryExecutor(ctrl)
	factory := deliveryapp.NewFactoryEventStrategy(executor)
	unknownErr := fmt.Errorf("some unknown err from factory")

	tests := []struct {
		name      string
		statusMsg string
		statusNow string
		wantType  any
		wantErr   error
	}{
		{
			"valid for assign strategy",
			deliveryapp.Created,
			"cooking",
			deliveryapp.AssignStrategy{executor},
			nil,
		},
		{
			"valid for unassign strategy",
			deliveryapp.Deleted,
			deliveryapp.Deleted,
			deliveryapp.UnassignStrategy{executor},
			nil,
		},
		{
			"valid for complete strategy",
			deliveryapp.Completed,
			deliveryapp.Completed,
			deliveryapp.CompleteStrategy{executor},
			nil,
		},
		{
			"unknown status",
			"pending",
			"",
			nil,
			unknownErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := factory.GetEventStrategy(tt.statusMsg, tt.statusNow)
			assert.IsType(t, tt.wantType, res)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
