package deliveryhttp_test

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"service-courier/internal/domain/delivery"
	"service-courier/internal/handler/deliveryhttp"
	"service-courier/observability/logger"
)

func TestDeliveryHandler_DeliveryAssign(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockdeliveryService(ctrl)
	zlog, err := logger.NewZapAdapter()
	if err != nil {
		log.Printf("failed to init logger: %v", err)
	}
	defer zlog.Sync()

	delivery := delivery.AssignResult{
		CourierID:     5,
		OrderID:       "8e6f9097-7c2e-4d84-ba28-0f3b5521a09c",
		TransportType: "scooter",
		Deadline:      func() time.Time { res, _ := time.Parse(time.RFC3339, "2025-08-06T13:15:00Z"); return res }(),
	}
	deliveryJSON := `
	{
		"courier_id": 5,
		"order_id": "8e6f9097-7c2e-4d84-ba28-0f3b5521a09c",
		"transport_type": "scooter",
		"delivery_deadline": "2025-08-06T13:15:00Z"
	}`

	tests := []struct {
		name         string
		inputBody    string
		expectedCode int
		expectedBody string
		behaviour    func(*MockdeliveryService)
	}{
		{
			"invalid json",
			`{""`,
			http.StatusBadRequest,
			`{"error": "order information is incorrect"}`,
			nil,
		},
		{
			"empty order id",
			`{"order_id": ""}`,
			http.StatusBadRequest,
			`{"error": "the order ID can't be empty"}`,
			nil,
		},
		{
			"invalid order id",
			`{"order_id": "123-abc"}`,
			http.StatusBadRequest,
			`{"error": "order ID is incorrect"}`,
			nil,
		},
		{
			"unknown error",
			`{"order_id": "8e6f9097-7c2e-4d84-ba28-0f3b5521a09c"}`,
			http.StatusInternalServerError,
			`{"error": "database error"}`,
			func(m *MockdeliveryService) {
				m.EXPECT().
					Assign(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("some unknown error from service"))
			},
		},
		{
			"valid",
			`{"order_id": "8e6f9097-7c2e-4d84-ba28-0f3b5521a09c"}`,
			http.StatusOK,
			deliveryJSON,
			func(m *MockdeliveryService) {
				m.EXPECT().
					Assign(gomock.Any(), gomock.Any()).
					Return(&delivery, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.behaviour != nil {
				tt.behaviour(m)
			}
			h := deliveryhttp.NewDeliveryHandler(zlog, m)

			r := httptest.NewRequest(http.MethodPost, "/delivery/assign", strings.NewReader(tt.inputBody))
			w := httptest.NewRecorder()

			h.Assign(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

		})
	}
}

func TestDeliveryHandler_DeliveryUnassign(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockdeliveryService(ctrl)
	zlog, err := logger.NewZapAdapter()
	if err != nil {
		log.Printf("failed to init logger: %v", err)
	}
	defer zlog.Sync()

	delivery := delivery.UnassignResult{
		CourierID: 5,
		OrderID:   "8e6f9097-7c2e-4d84-ba28-0f3b5521a09c",
		Status:    "busy",
	}
	deliveryJSON := `
	{
		"courier_id": 5,
		"order_id": "8e6f9097-7c2e-4d84-ba28-0f3b5521a09c",
		"status": "busy"
	}`

	tests := []struct {
		name         string
		inputBody    string
		expectedCode int
		expectedBody string
		behaviour    func(*MockdeliveryService)
	}{
		{
			"invalid json",
			`{""`,
			http.StatusBadRequest,
			`{"error": "order information is incorrect"}`,
			nil,
		},
		{
			"empty order id",
			`{"order_id": ""}`,
			http.StatusBadRequest,
			`{"error": "the order ID can't be empty"}`,
			nil,
		},
		{
			"invalid order id",
			`{"order_id": "123-abc"}`,
			http.StatusBadRequest,
			`{"error": "order ID is incorrect"}`,
			nil,
		},
		{
			"unknown error",
			`{"order_id": "8e6f9097-7c2e-4d84-ba28-0f3b5521a09c"}`,
			http.StatusInternalServerError,
			`{"error": "database error"}`,
			func(m *MockdeliveryService) {
				m.EXPECT().
					Unassign(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("some unknown error from service"))
			},
		},
		{
			"valid",
			`{"order_id": "8e6f9097-7c2e-4d84-ba28-0f3b5521a09c"}`,
			http.StatusOK,
			deliveryJSON,
			func(m *MockdeliveryService) {
				m.EXPECT().
					Unassign(gomock.Any(), gomock.Any()).
					Return(&delivery, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.behaviour != nil {
				tt.behaviour(m)
			}
			h := deliveryhttp.NewDeliveryHandler(zlog, m)

			r := httptest.NewRequest(http.MethodPost, "/delivery/assign", strings.NewReader(tt.inputBody))
			w := httptest.NewRecorder()

			h.Unassign(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

		})
	}
}
