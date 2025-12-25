package courierhttp_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"service-courier/internal/domain/courier"
	"service-courier/internal/handler/courierhttp"
	"service-courier/observability/logger"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCourierHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockcourierService(ctrl)
	log, _ := logger.NewZapAdapter()

	tests := []struct {
		name         string
		inputBody    string
		expectedCode int
		expectedBody string
		behaviour    func(*MockcourierService)
	}{
		{
			"invalid json",
			`{""`,
			http.StatusBadRequest,
			`{"error": "courier's information is incorrect"}`,
			nil,
		},
		{
			"empty name",
			`{"name": ""}`,
			http.StatusBadRequest,
			`{"error": "required courier fields aren't filled"}`,
			nil,
		},
		{
			"empty phone",
			`{"phone": ""}`,
			http.StatusBadRequest,
			`{"error": "required courier fields aren't filled"}`,
			nil,
		},
		{
			"invalid phone",
			`{"phone": "abc", "name": "TestName"}`,
			http.StatusBadRequest,
			`{"error": "courier's phone number is incorrect"}`,
			nil,
		},
		{
			"phone exist",
			`{"name": "TestName", "phone": "+1234567890"}`,
			http.StatusConflict,
			`{"error": "courier with such a phone already exists"}`,
			func(m *MockcourierService) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(int64(-1), courier.ErrCourierExistPhone)
			},
		},
		{
			"unknown error",
			`{"name": "TestName", "phone": "+1234567890"}`,
			http.StatusInternalServerError,
			`{"error": "database error"}`,
			func(m *MockcourierService) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(int64(-1), errors.New("some unknown error from service"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.behaviour != nil {
				tt.behaviour(m)
			}
			h := courierhttp.NewCourierHandler(log, m)

			r := httptest.NewRequest(http.MethodPost, "/courier", strings.NewReader(tt.inputBody))
			w := httptest.NewRecorder()

			h.Create(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

		})
	}
}

func TestCourierHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockcourierService(ctrl)
	log, _ := logger.NewZapAdapter()

	tests := []struct {
		name         string
		inputBody    string
		expectedCode int
		expectedBody string
		behaviour    func(*MockcourierService)
	}{
		{
			"invalid json",
			`{""`,
			http.StatusBadRequest,
			`{"error": "courier's information is incorrect"}`,
			nil,
		},
		{
			"invalid id",
			`{"id": -1}`,
			http.StatusBadRequest,
			`{"error": "courier's ID is incorrect"}`,
			nil,
		},
		{
			"empty id",
			`{"id": 0}`,
			http.StatusBadRequest,
			`{"error": "courier's ID is incorrect"}`,
			nil,
		},
		{
			"invalid phone",
			`{"id": 1, "phone": "abc", "name": "TestName"}`,
			http.StatusBadRequest,
			`{"error": "courier's phone number is incorrect"}`,
			nil,
		},
		{
			"phone exist",
			`{"id": 1, "name": "TestName", "phone": "+1234567890"}`,
			http.StatusConflict,
			`{"error": "courier with such a phone already exists"}`,
			func(m *MockcourierService) {
				m.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(int64(-1), courier.ErrCourierExistPhone)
			},
		},
		{
			"courier for update not found",
			`{"id": 100}`,
			http.StatusNotFound,
			`{"error": "courier wasn't found"}`,
			func(m *MockcourierService) {
				m.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(int64(-1), courier.ErrCourierNotFound)
			},
		},
		{
			"unknown error",
			`{"id": 1, "name": "TestName", "phone": "+1234567890"}`,
			http.StatusInternalServerError,
			`{"error": "database error"}`,
			func(m *MockcourierService) {
				m.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(int64(-1), errors.New("some unknown error from service"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.behaviour != nil {
				tt.behaviour(m)
			}
			h := courierhttp.NewCourierHandler(log, m)

			r := httptest.NewRequest(http.MethodPut, "/courier", strings.NewReader(tt.inputBody))
			w := httptest.NewRecorder()

			h.Update(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

		})
	}
}

func TestCourierHandler_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockcourierService(ctrl)
	log, _ := logger.NewZapAdapter()

	c := courier.Courier{
		ID:            1,
		Name:          "TestName",
		Phone:         "+1234567890",
		Status:        "busy",
		TransportType: "scooter",
	}
	courierJSON := `
	{
		"id": 1,
		"name": "TestName",
		"phone": "+1234567890",
		"status": "busy",
		"transport_type": "scooter"
	}`

	tests := []struct {
		name         string
		inputID      string
		expectedCode int
		expectedBody string
		behaviour    func(*MockcourierService)
	}{
		{
			"invalid id",
			"-1",
			http.StatusBadRequest,
			`{"error": "courier's ID is incorrect"}`,
			nil,
		},
		{
			"empty id",
			"",
			http.StatusBadRequest,
			`{"error": "courier's ID is incorrect"}`,
			nil,
		},
		{
			"courier not found",
			"100",
			http.StatusNotFound,
			`{"error": "courier wasn't found"}`,
			func(m *MockcourierService) {
				m.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(nil, courier.ErrCourierNotFound)
			},
		},
		{
			"unknown error",
			"1",
			http.StatusInternalServerError,
			`{"error": "database error"}`,
			func(m *MockcourierService) {
				m.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("some unknown error from service"))
			},
		},
		{
			"valid",
			"1",
			http.StatusOK,
			courierJSON,
			func(m *MockcourierService) {
				m.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(&c, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.behaviour != nil {
				tt.behaviour(m)
			}
			h := courierhttp.NewCourierHandler(log, m)

			r := httptest.NewRequest(http.MethodGet, "/courier/1", nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.inputID)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			h.GetByID(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

		})
	}
}

func TestCourierHandler_GetMulti(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockcourierService(ctrl)
	log, _ := logger.NewZapAdapter()

	tests := []struct {
		name         string
		expectedCode int
		expectedBody string
		behaviour    func(*MockcourierService)
	}{
		{
			"unknown error",
			http.StatusInternalServerError,
			`{"error": "database error"}`,
			func(m *MockcourierService) {
				m.EXPECT().
					GetMulti(gomock.Any()).
					Return(nil, errors.New("some unknown error from service"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.behaviour != nil {
				tt.behaviour(m)
			}
			h := courierhttp.NewCourierHandler(log, m)

			r := httptest.NewRequest(http.MethodGet, "/couriers", nil)
			w := httptest.NewRecorder()

			h.GetMulti(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

		})
	}
}
