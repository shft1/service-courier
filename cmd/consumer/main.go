package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"service-courier/internal/config/consumercfg"
	"service-courier/internal/config/dbcfg"
	"service-courier/internal/databus/kafka"
	"service-courier/internal/db/postgre"
	"service-courier/internal/gateway/ordergrpc"
	"service-courier/internal/handler/orderbus"
	"service-courier/internal/proto/orderpb"
	"service-courier/internal/repository/courierdb"
	"service-courier/internal/repository/deliverydb"
	"service-courier/internal/service/deliveryapp"
	"service-courier/observability/logger"
	"syscall"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	pool := postgre.InitPool(sysCtx, zlog, dbEnv)
	defer pool.Close()

	// Инициализация менеджера транзакций
	txManager := postgre.NewTxManagerPostgre(pool)

	// Инициализация фабрики времени
	timeFactory := deliveryapp.NewFactoryTimeCalculator()

	// Инициализация репозитория курьера
	courDB := courierdb.NewCourierRepository(pool, txManager)

	// Инициализация сервиса доставок
	delDB := deliverydb.NewDeliveryRepository(pool, txManager)
	delApp := deliveryapp.NewDeliveryService(delDB, courDB, timeFactory, txManager)

	// Инициализация фабрики бизнес-операций
	eventFactory := deliveryapp.NewFactoryEventStrategy(delApp)

	// Инициализация gRPC соединения
	grpcServer := fmt.Sprintf("%v:%v", consumEnv.OrderHost, consumEnv.OrderPort)
	conn, err := grpc.NewClient(grpcServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	kafkaClient, err := kafka.NewKafkaClient(zlog, consumEnv, handler, []string{consumEnv.KafkaTopic})
	if err != nil {
		zlog.Error("failed to create Kafka client", logger.NewField("error", err))
		return
	}
	defer kafkaClient.Close()

	// Запуск консьюминга Kafka
	go kafkaClient.Consume(sysCtx)

	zlog.Info("start kafka consuming")

	// Ожидание отмены контекста
	<-sysCtx.Done()

	zlog.Info("stopped kafka consuming")
}
