//go:build integration

package integration_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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
	"service-courier/internal/handler/courierhttp"
	"service-courier/internal/repository/courierdb"
	"service-courier/internal/router/courierroute"
	"service-courier/internal/service/courierapp"
)

const (
	createURL  = "/courier"
	getByIDURL = "/courier/"
)

type CourierTestSuite struct {
	suite.Suite
	server    *httptest.Server
	pool      *pgxpool.Pool
	ctx       context.Context
	courierID int64
}

func (s *CourierTestSuite) SetupSuite() {
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

	repo := courierdb.NewCourierRepository(pool, txManager)
	srv := courierapp.NewCourierService(repo)
	hand := courierhttp.NewCourierHandler(srv)
	router := chi.NewRouter()
	courierroute.CourierRoute(router, hand)

	s.server = httptest.NewServer(router)
}

func (s *CourierTestSuite) SetupTest() {
	err := s.pool.QueryRow(
		s.ctx,
		`INSERT INTO couriers (name, phone, status, transport_type)
			VALUES ($1, $2, $3, $4)
			RETURNING id;`,
		"TestCourier", "+1234567890", "busy", "scooter",
	).Scan(&s.courierID)
	s.Require().NoError(err)
}

func (s *CourierTestSuite) TestCreateCourier() {
	payload := `{"name": "TestName", "phone": "+0234567890"}`
	resp, err := http.Post(s.server.URL+createURL, "application/json", strings.NewReader(payload))
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.EqualValues(0, resp.ContentLength)
}

func (s *CourierTestSuite) TestGetByIDCourier() {
	resp, err := http.Get(s.server.URL + getByIDURL + strconv.Itoa(int(s.courierID)))
	s.Require().NoError(err)
	defer resp.Body.Close()

	var cour courierhttp.CourierResponse
	err = json.NewDecoder(resp.Body).Decode(&cour)
	s.Require().NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(s.courierID, cour.ID)
	s.Equal("TestCourier", cour.Name)
	s.Equal("+1234567890", cour.Phone)
	s.Equal("busy", cour.Status)
	s.Equal("scooter", cour.TransportType)
}

func (s *CourierTestSuite) TearDownTest() {
	_, err := s.pool.Exec(
		s.ctx,
		`DELETE FROM couriers WHERE id = $1`,
		s.courierID,
	)
	s.Require().NoError(err)
}

func (s *CourierTestSuite) TearDownSuite() {
	s.server.Close()
}

func TestCourier(t *testing.T) {
	suite.Run(t, new(CourierTestSuite))
}
