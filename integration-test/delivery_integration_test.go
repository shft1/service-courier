//go:build integration

package integration_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"service-courier/internal/db/postgre"
	"service-courier/internal/entity/delivery"
	deliveryHandler "service-courier/internal/handler/delivery"
	courierRepository "service-courier/internal/repository/courier"
	deliveryRepository "service-courier/internal/repository/delivery"
	deliveryRouter "service-courier/internal/router/delivery"
	deliveryService "service-courier/internal/service/delivery"
)

const (
	assignURL   = "/delivery/assign"
	unassignURL = "/delivery/unassign"
)

type DeliveryTestSuite struct {
	suite.Suite
	server        *httptest.Server
	pool          *pgxpool.Pool
	ctx           context.Context
	orderID       string
	courierID     int
	freeCourierID int
	newDeliveryID string
}

func (s *DeliveryTestSuite) SetupSuite() {
	s.ctx = context.Background()

	pg, err := postgres.Run(s.ctx,
		"postgres:15",
		postgres.WithDatabase("testDB"),
		postgres.WithUsername("testUser"),
		postgres.WithPassword("testPass"),
	)
	s.Require().NoError(err)

	s.T().Cleanup(func() {
		pg.Terminate(s.ctx)
		fmt.Println("info: test-containter terminated")
	})

	dsn, err := pg.ConnectionString(s.ctx)
	s.Require().NoError(err)

	cfg, err := pgxpool.ParseConfig(dsn)
	s.Require().NoError(err)
	cfg.MaxConns = 1

	time.Sleep(1 * time.Second)

	pool, err := pgxpool.NewWithConfig(s.ctx, cfg)
	s.Require().NoError(err)
	s.pool = pool
	s.T().Cleanup(func() {
		pool.Close()
		fmt.Println("info: pgxpool closed")
	})

	stdpool, err := sql.Open("pgx", dsn)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		stdpool.Close()
		fmt.Println("info: stdpool closed")
	})

	err = goose.Up(stdpool, "../migrations")
	s.Require().NoError(err)

	txManager := postgre.NewTxManagerPostgre(pool)

	courierRepository := courierRepository.NewCourierRepository(pool, txManager)

	deliveryRepository := deliveryRepository.NewDeliveryRepository(pool, txManager)
	deliveryService := deliveryService.NewDeliveryService(deliveryRepository, courierRepository, txManager)
	deliveryHandler := deliveryHandler.NewDeliveryHandler(deliveryService)

	router := chi.NewRouter()
	deliveryRouter.DeliveryRoute(router, deliveryHandler)
	s.server = httptest.NewServer(router)
}

func (s *DeliveryTestSuite) SetupTest() {
	tx, err := s.pool.Begin(s.ctx)
	s.Require().NoError(err)

	err = tx.QueryRow(
		s.ctx,
		`INSERT INTO couriers (name, phone, status, transport_type) VALUES ($1, $2, $3, $4) RETURNING id;`,
		"TestCourier2", "+01234567890", "available", "car",
	).Scan(&s.freeCourierID)
	s.Require().NoError(err)

	err = tx.QueryRow(
		s.ctx,
		`INSERT INTO couriers (name, phone, status, transport_type) VALUES ($1, $2, $3, $4) RETURNING id;`,
		"TestCourier", "+1234567890", "busy", "scooter",
	).Scan(&s.courierID)
	s.Require().NoError(err)

	err = tx.QueryRow(
		s.ctx,
		`INSERT INTO delivery (courier_id, order_id, deadline) VALUES ($1, $2, now()) RETURNING courier_id, order_id;`,
		s.courierID, "8e6f9097-7c2e-4d84-ba28-0f3b5521a09c",
	).Scan(&s.courierID, &s.orderID)
	s.Require().NoError(err)

	err = tx.Commit(s.ctx)
	s.Require().NoError(err)
}

func (s *DeliveryTestSuite) TestDeliveryAssign() {
	payload := `{"order_id": "1e4f9095-7c2e-4d84-ba28-0f3b5521a19c"}`
	resp, err := http.Post(s.server.URL+assignURL, "application/json", strings.NewReader(payload))
	s.Require().NoError(err)

	var delivery delivery.DeliveryAssign
	err = json.NewDecoder(resp.Body).Decode(&delivery)
	s.Require().NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal("1e4f9095-7c2e-4d84-ba28-0f3b5521a19c", delivery.OrderID)
	s.Equal(s.freeCourierID, delivery.CourierID)
	s.NotZero(delivery.TransportType)
	s.NotZero(delivery.DeliveryDeadline)

	s.newDeliveryID = delivery.OrderID
}

func (s *DeliveryTestSuite) TestDeliveryUnassign() {
	payload := fmt.Sprintf(`{"order_id": "%s"}`, s.orderID)
	resp, err := http.Post(s.server.URL+unassignURL, "application/json", strings.NewReader(payload))
	s.Require().NoError(err)

	var delivery delivery.DeliveryUnassign
	err = json.NewDecoder(resp.Body).Decode(&delivery)
	s.Require().NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(s.orderID, delivery.OrderID)
	s.Equal("unassigned", delivery.Status)
	s.Equal(s.courierID, delivery.CourierID)
}

func (s *DeliveryTestSuite) TearDownTest() {
	tx, err := s.pool.Begin(s.ctx)
	s.Require().NoError(err)

	_, err = tx.Exec(
		s.ctx,
		`DELETE FROM delivery WHERE order_id = $1`,
		s.newDeliveryID,
	)
	s.Require().NoError(err)

	_, err = tx.Exec(
		s.ctx,
		`DELETE FROM couriers WHERE id = $1`,
		s.freeCourierID,
	)
	s.Require().NoError(err)

	_, err = tx.Exec(
		s.ctx,
		`DELETE FROM delivery WHERE order_id = $1`,
		s.orderID,
	)
	s.Require().NoError(err)

	_, err = tx.Exec(
		s.ctx,
		`DELETE FROM couriers WHERE id = $1`,
		s.courierID,
	)
	s.Require().NoError(err)

	err = tx.Commit(s.ctx)
	s.Require().NoError(err)
}

func (s *DeliveryTestSuite) TearDownSuite() {
	s.server.Close()
}

func TestDelivery(t *testing.T) {
	suite.Run(t, new(DeliveryTestSuite))
}
