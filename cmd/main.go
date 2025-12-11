package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"service-courier/internal/cli"
	"service-courier/internal/config"
	"service-courier/internal/db/postgre"
	"service-courier/internal/gateway/order"
	"service-courier/internal/handler/courierhttp"
	"service-courier/internal/handler/deliveryhttp"
	"service-courier/internal/handler/healthhttp"
	orderPB "service-courier/internal/proto/order"
	"service-courier/internal/repository/courierdb"
	"service-courier/internal/repository/deliverydb"
	"service-courier/internal/router"
	"service-courier/internal/server"
	"service-courier/internal/service/courierapp"
	"service-courier/internal/service/deliveryapp"
	"service-courier/internal/worker"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Инициализация основного контекста
	sysCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Инициализация env переменных
	env := config.SetupEnv()

	// Инициализация пула соединений с БД
	pool := postgre.InitPool(sysCtx, env)
	defer pool.Close()

	// Инициализация менеджера транзакций
	txManager := postgre.NewTxManagerPostgre(pool)

	// Инициализация логики состояния
	healthHTTP := healthhttp.NewHealthHandler()

	// Инициализация логики курьеров
	courDB := courierdb.NewCourierRepository(pool, txManager)
	courApp := courierapp.NewCourierService(courDB)
	courHTTP := courierhttp.NewCourierHandler(courApp)

	// Инициализация фабрики
	timeFactory := deliveryapp.NewFactoryTimeCalculator()

	// Инициализация логики доставок
	delDB := deliverydb.NewDeliveryRepository(pool, txManager)
	delApp := deliveryapp.NewDeliveryService(delDB, courDB, timeFactory, txManager)
	delHTTP := deliveryhttp.NewDeliveryHandler(delApp)

	// Регистрация адресов на обработчиков
	router := router.SetupRoute(healthHTTP, courHTTP, delHTTP)

	// Инициализация gRPC
	conn, err := grpc.NewClient(env.OrderPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("failed to connect to gRPC server", err)
	}
	defer conn.Close()

	clientPB := orderPB.NewOrdersServiceClient(conn)
	orderGW := order.NewGateway(clientPB)

	// Инициализация фонового воркера получения и назначения заказов
	pollPeriod, err := time.ParseDuration(env.TimePoll)
	if err != nil {
		fmt.Println("error parsing time poll period:", err)
		pollPeriod = time.Second * 10
	}
	orderPoller := worker.NewOrderPoller(pollPeriod, orderGW, delApp)

	go orderPoller.Start(sysCtx)

	// Инициализация фонового воркера проверки доставок
	checkPeriod, err := time.ParseDuration(env.TimeCheck)
	if err != nil {
		fmt.Println("error parsing time check duration:", err)
		checkPeriod = time.Second * 10
	}
	deliveryChecker := worker.NewDeliveryMonitor(checkPeriod, delApp)

	go deliveryChecker.Start(sysCtx)

	// Парсинг командной строки
	cmd := cli.CliHandler(env)
	if err := cmd.Run(sysCtx, os.Args); err != nil {
		fmt.Println(err)
	}

	// Запуск сервера через graceful shutdown
	server.StartServerGraceful(sysCtx, router, pool, env)
}
