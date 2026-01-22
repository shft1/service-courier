package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"service-courier/internal/config/consumercfg"
	"service-courier/internal/config/dbcfg"
	"service-courier/internal/databus/kafka"
	"service-courier/internal/db/postgre"
	"service-courier/internal/gateway/ordergrpc"
	"service-courier/internal/handler/orderbus"
	"service-courier/internal/middleware/mdrpc"
	"service-courier/internal/proto/orderpb"
	"service-courier/internal/repository/courierdb"
	"service-courier/internal/repository/deliverydb"
	"service-courier/internal/resilience/retry"
	"service-courier/internal/router/metricsroute"
	"service-courier/internal/server"
	"service-courier/internal/service/deliveryapp"
	"service-courier/observability/logger"
	"service-courier/observability/metrics/metricsrpc"
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
	// Инициализация env переменных базы данных
	dbEnv := dbcfg.SetupDataBaseEnv()

	// Инициализация env переменных консьюмера
	consumEnv := consumercfg.SetupConsumerEnv(zlog)

	// Инициализация пула соединений с БД
	pool, err := postgre.InitPool(sysCtx, zlog, dbEnv)
	if err != nil {
		zlog.Error("failed to create connection pool", logger.NewField("error", err))
		return
	}
	defer pool.Close()

	// Инициализация менеджера транзакций
	txManager := postgre.NewTxManagerPostgre(zlog, pool)

	// Инициализация фабрики времени
	timeFactory := deliveryapp.NewFactoryTimeCalculator()

	// Инициализация репозитория курьера
	courDB := courierdb.NewCourierRepository(pool, txManager)

	// Инициализация сервиса доставок
	delDB := deliverydb.NewDeliveryRepository(pool, txManager)
	delApp := deliveryapp.NewDeliveryService(deliveryapp.Arguments{
		DelRepo: delDB, CourRepo: courDB, Factory: timeFactory, TxManager: txManager,
	})

	// Инициализация фабрики бизнес-операций
	eventFactory := deliveryapp.NewFactoryEventStrategy(delApp)

	// Инициализация Logger интерцептора
	loggerInter := mdrpc.NewLoggerInterceptor(zlog)

	// Инициализация Retry интерцептора
	strategy := retry.NewExponentialBackoffWithJitter(retry.Arguments{
		Multi:     consumEnv.Multiplier,
		Jitter:    consumEnv.Jitter,
		InitDelay: consumEnv.InitDelay,
		MaxDelay:  consumEnv.MaxDelay,
	})
	retry := retry.NewRetryExecutor(
		retry.WithMaxAttempts(consumEnv.MaxAttempts),
		retry.WithStrategy(strategy),
		retry.WithShouldRetry(retry.ShouldRetry),
	)
	retryInter := mdrpc.NewRetryInterceptor(retry)

	// Инициализация Metrics RPC
	metrics := metricsrpc.NewRPCMetrics()

	// Инициализация Metrics интерцептора
	metricsInter := mdrpc.NewMetricsInterceptor(metrics, retry.IsRetryFromContext)

	// Инициализация gRPC соединения
	grpcServer := net.JoinHostPort(consumEnv.OrderHost, consumEnv.OrderPort)
	conn, err := grpc.NewClient(
		grpcServer,
		grpc.WithChainUnaryInterceptor(loggerInter, retryInter, metricsInter),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		zlog.Error("failed to connect to gRPC server", logger.NewField("error", err))
		return
	}
	defer conn.Close()

	// Инициализация gRPC клиента
	clientPB := orderpb.NewOrdersServiceClient(conn)
	orderGW := ordergrpc.NewGateway(clientPB)

	// Инициализация обработчика топика changed Kafka
	handler := orderbus.NewConsumeHandler(zlog, orderGW, eventFactory)

	// Инициализация Kafka клиента
	kafkaClient, err := kafka.NewKafkaClient(kafka.Arguments{
		Log: zlog, Env: consumEnv, Handler: handler, Topics: []string{consumEnv.KafkaTopic},
	})
	if err != nil {
		zlog.Error("failed to create Kafka client", logger.NewField("error", err))
		return
	}
	defer kafkaClient.Close()

	// Запуск консьюминга Kafka
	go kafkaClient.Consume(sysCtx)
	zlog.Info("start kafka consuming")

	// Инициализация роутера
	router := chi.NewRouter()

	// Инициализация обработчика метрик
	metricsHTTP := promhttp.Handler().ServeHTTP

	// Регистрация обработчика метрик в роутере
	metricsroute.MetricsRoute(router, metricsHTTP)

	// Запуск сервера через graceful shutdown
	server.StartServerGraceful(sysCtx, zlog, router, consumEnv.Port)
}
