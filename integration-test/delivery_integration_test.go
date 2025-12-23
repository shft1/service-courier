//go:build integration

package integration_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
	"service-courier/internal/domain/delivery"
	"service-courier/internal/handler/deliveryhttp"
	"service-courier/internal/repository/courierdb"
	"service-courier/internal/repository/deliverydb"
	"service-courier/internal/router/deliveryroute"
	"service-courier/internal/service/deliveryapp"
	"service-courier/observability/logger"
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
	courierID     int64
	freeCourierID int64
	newDeliveryID string
}

func (s *DeliveryTestSuite) SetupSuite() {
	s.ctx = context.Background()

	zlog, err := logger.NewZapAdapter()
	if err != nil {
		log.Printf("failed to init logger: %v", err)
	}
	defer zlog.Sync()

	pg, err := postgres.Run(s.ctx,
		"postgres:15",
		postgres.WithDatabase("testDB"),
		postgres.WithUsername("testUser"),
		postgres.WithPassword("testPass"),
	)
	s.Require().NoError(err)

	s.T().Cleanup(func() {
		pg.Terminate(s.ctx)
		zlog.Info("test-containter terminated")
	})

	dsn, err := pg.ConnectionString(s.ctx)
	s.Require().NoError(err)

	cfg, err := pgxpool.ParseConfig(dsn)
	s.Require().NoError(err)
	cfg.MaxConns = 1

	time.Sleep(2 * time.Second)

	pool, err := pgxpool.NewWithConfig(s.ctx, cfg)
	s.Require().NoError(err)
	s.pool = pool
	s.T().Cleanup(func() {
		pool.Close()
		zlog.Info("test-pgxpool closed")
	})

	stdpool, err := sql.Open("pgx", dsn)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		stdpool.Close()
		zlog.Info("test-stdpool closed")
	})

	err = goose.Up(stdpool, "../migrations")
	s.Require().NoError(err)

	txManager := postgre.NewTxManagerPostgre(pool)

	timeFactory := deliveryapp.NewFactoryTimeCalculator()

	courierdb := courierdb.NewCourierRepository(pool, txManager)

	deliverydb := deliverydb.NewDeliveryRepository(pool, txManager)
	deliveryapp := deliveryapp.NewDeliveryService(deliverydb, courierdb, timeFactory, txManager)
	deliveryhttp := deliveryhttp.NewDeliveryHandler(deliveryapp)

	router := chi.NewRouter()
	deliveryroute.DeliveryRoute(router, deliveryhttp)
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
	defer resp.Body.Close()

	var del deliveryhttp.DeliveryAssignResponse
	err = json.NewDecoder(resp.Body).Decode(&del)
	s.Require().NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal("1e4f9095-7c2e-4d84-ba28-0f3b5521a19c", del.OrderID)
	s.Equal(s.freeCourierID, del.CourierID)
	s.NotZero(del.TransportType)
	s.NotZero(del.Deadline)

	s.newDeliveryID = del.OrderID
}

func (s *DeliveryTestSuite) TestDeliveryUnassign() {
	payload := fmt.Sprintf(`{"order_id": "%s"}`, s.orderID)
	resp, err := http.Post(s.server.URL+unassignURL, "application/json", strings.NewReader(payload))
	s.Require().NoError(err)
	defer resp.Body.Close()

	var del deliveryhttp.DeliveryUnassignResponse
	err = json.NewDecoder(resp.Body).Decode(&del)
	s.Require().NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(s.orderID, del.OrderID)
	s.Equal(delivery.UnassignStatus, del.Status)
	s.Equal(s.courierID, del.CourierID)
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
