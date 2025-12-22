package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"service-courier/internal/cli"
	"service-courier/internal/config/appcfg"
	"service-courier/internal/config/dbcfg"
	"service-courier/internal/db/postgre"
	"service-courier/internal/handler/courierhttp"
	"service-courier/internal/handler/deliveryhttp"
	"service-courier/internal/handler/healthhttp"
	"service-courier/internal/repository/courierdb"
	"service-courier/internal/repository/deliverydb"
	"service-courier/internal/router"
	"service-courier/internal/router/middleware"
	"service-courier/internal/server"
	"service-courier/internal/service/courierapp"
	"service-courier/internal/service/deliveryapp"
	"service-courier/internal/worker/deliveryworker"
	"service-courier/observability/logger"
	"service-courier/observability/metrics"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	pool := postgre.InitPool(sysCtx, zlog, dbEnv)
	defer pool.Close()

	// Инициализация менеджера транзакций
	txManager := postgre.NewTxManagerPostgre(pool)

	// Инициализация логики состояния
	healthHTTP := healthhttp.NewHealthHandler()

	// Инициализация логики курьеров
	courDB := courierdb.NewCourierRepository(pool, txManager)
	courApp := courierapp.NewCourierService(courDB)
	courHTTP := courierhttp.NewCourierHandler(courApp)

	// Инициализация фабрики времени
	timeFactory := deliveryapp.NewFactoryTimeCalculator()

	// Инициализация логики доставок
	delDB := deliverydb.NewDeliveryRepository(pool, txManager)
	delApp := deliveryapp.NewDeliveryService(delDB, courDB, timeFactory, txManager)
	delHTTP := deliveryhttp.NewDeliveryHandler(delApp)

	// Инициализация Middleware логгирования
	loggerMW := middleware.NewLoggerMiddleware(zlog)

	// Инициализация метрик
	metrics := metrics.NewHTTPMetrics()

	// Инициализация Middleware метрик
	metricsMW := middleware.NewMetricsMiddleware(metrics)

	// Инициализация обработчика метрик
	metricsHTTP := promhttp.Handler().ServeHTTP

	// Регистрация адресов и middleware
	router := router.SetupRoute(loggerMW, metricsMW, healthHTTP, courHTTP, delHTTP, metricsHTTP)

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
	go deliveryChecker.Start(sysCtx)

	// Парсинг командной строки
	cmd := cli.CliHandler(appEnv)
	if err := cmd.Run(sysCtx, os.Args); err != nil {
		zlog.Error("failed to parse cli command", logger.NewField("error", err))
	}
	// Запуск сервера через graceful shutdown
	server.StartServerGraceful(zlog, sysCtx, router, pool, appEnv)
}
