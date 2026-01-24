package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"service-courier/internal/cli"
	"service-courier/internal/config/appcfg"
	"service-courier/internal/config/dbcfg"
	"service-courier/internal/db/postgre"
	"service-courier/internal/handler/courierhttp"
	"service-courier/internal/handler/deliveryhttp"
	"service-courier/internal/handler/healthhttp"
	"service-courier/internal/middleware/mdhttp"
	"service-courier/internal/repository/courierdb"
	"service-courier/internal/repository/deliverydb"
	"service-courier/internal/resilience/limiter"
	"service-courier/internal/router"
	"service-courier/internal/router/pprofroute"
	"service-courier/internal/server"
	"service-courier/internal/service/courierapp"
	"service-courier/internal/service/deliveryapp"
	"service-courier/internal/worker/deliveryworker"
	"service-courier/observability/logger"
	"service-courier/observability/metrics/metricshttp"
)

func main() {
	// Инициализация основного контекста
	sysCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Инициализация логгера
	zlog, err := logger.NewZapAdapter()
	if err != nil {
		log.Printf("failed to init logger: %v", err)
	}
	defer zlog.Sync()

	// Загрузка env переменных в окружение
	if err := godotenv.Load(); err != nil {
		zlog.Error("Error loading .env file")
	}
	// Инициализация env переменных приложения
	appEnv := appcfg.SetupAppEnv()

	// Инициализация env переменных базы данных
	dbEnv := dbcfg.SetupDataBaseEnv()

	// Инициализация пула соединений с БД
	pool, err := postgre.InitPool(sysCtx, zlog, dbEnv)
	if err != nil {
		zlog.Error("failed to create connection pool", logger.NewField("error", err))
		return
	}
	defer pool.Close()

	// Инициализация менеджера транзакций
	txManager := postgre.NewTxManagerPostgre(zlog, pool)

	// Инициализация логики состояния
	healthHTTP := healthhttp.NewHealthHandler(zlog)

	// Инициализация логики курьеров
	courDB := courierdb.NewCourierRepository(pool, txManager)
	courApp := courierapp.NewCourierService(courDB)
	courHTTP := courierhttp.NewCourierHandler(zlog, courApp)

	// Инициализация фабрики времени
	timeFactory := deliveryapp.NewFactoryTimeCalculator()

	// Инициализация логики доставок
	delDB := deliverydb.NewDeliveryRepository(pool, txManager)
	delApp := deliveryapp.NewDeliveryService(deliveryapp.Arguments{
		DelRepo: delDB, CourRepo: courDB, Factory: timeFactory, TxManager: txManager,
	})
	delHTTP := deliveryhttp.NewDeliveryHandler(zlog, delApp)

	// Инициализация Middleware логгирования
	loggerMW := mdhttp.NewLoggerMiddleware(zlog)

	// Инициализация метрик
	metricshttp := metricshttp.NewHTTPMetrics()

	// Инициализация Middleware метрик
	metricsMW := mdhttp.NewMetricsMiddleware(metricshttp)

	// Инициализация обработчика метрик
	metricsHTTP := promhttp.Handler().ServeHTTP

	refil, err := time.ParseDuration(appEnv.Refill)
	if err != nil {
		zlog.Error("failed to parse refil time", logger.NewField("error", err))
		return
	}
	limit, err := strconv.Atoi(appEnv.Limit)
	if err != nil {
		zlog.Error("failed to parse limit bucket", logger.NewField("error", err))
		return
	}

	// Инициализация Middleware ограничителя запросов
	limiter := limiter.NewTokenBucketLimiter(refil, limit)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		limiter.StartReplenishment(sysCtx)
	}()

	limiterMW := mdhttp.NewLimiterMiddleware(zlog, limiter)

	// Регистрация адресов и middleware
	router := router.SetupRoute(loggerMW, metricsMW, limiterMW, healthHTTP, courHTTP, delHTTP, metricsHTTP)

	// // [Note] - Работа с заказами происходит через Kafka
	// // Инициализация gRPC соединения
	// conn, err := grpc.NewClient(appEnv.OrderPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	zlog.Error("failed to connect to gRPC server", logger.NewField("error", err))
	// }
	// defer conn.Close()

	// // Инициализация gRPC клиента
	// clientPB := orderpb.NewOrdersServiceClient(conn)
	// orderGW := ordergrpc.NewGateway(clientPB)

	// // Инициализация фонового воркера получения и назначения заказов
	// pollPeriod, err := time.ParseDuration(appEnv.TimePoll)
	// if err != nil {
	// 	zlog.Warn("failed to parse duration", logger.NewField("error", err))
	// 	pollPeriod = time.Second * 10
	// }
	// orderPoller := orderworker.NewOrderPoller(zlog, pollPeriod, orderGW, delApp)

	// go orderPoller.Start(sysCtx)

	// Инициализация фонового воркера проверки доставок
	checkPeriod, err := time.ParseDuration(appEnv.TimeCheck)
	if err != nil {
		zlog.Warn("failed to parse duration", logger.NewField("error", err))
		checkPeriod = time.Second * 10
	}
	deliveryChecker := deliveryworker.NewDeliveryMonitor(zlog, checkPeriod, delApp)

	// Запуск фоновой проверки доставок
	wg.Add(1)
	go func() {
		defer wg.Done()
		deliveryChecker.Start(sysCtx)
	}()

	// Парсинг командной строки
	cmd := cli.CliHandler(appEnv)
	if err := cmd.Run(sysCtx, os.Args); err != nil {
		zlog.Error("failed to parse cli command", logger.NewField("error", err))
		return
	}

	// Запуск основного веб-сервера
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.StartServer(sysCtx, zlog, router, appEnv.AppHost, appEnv.AppPort)
	}()

	// Инициализация роутера и регистрация pprof-адресов
	proute := chi.NewRouter()
	pprofroute.PprofRoute(proute)

	// Запуск pprof веб-сервера
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.StartPprofServer(sysCtx, zlog, proute)
	}()

	wg.Wait()
	zlog.Info("service-courier has been stopted")
}
