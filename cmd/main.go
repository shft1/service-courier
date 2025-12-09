package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"service-courier/internal/cli"
	"service-courier/internal/config"
	"service-courier/internal/db/postgre"
	"service-courier/internal/handler/courierhttp"
	"service-courier/internal/handler/deliveryhttp"
	"service-courier/internal/handler/healthhttp"
	"service-courier/internal/repository/courierdb"
	"service-courier/internal/repository/deliverydb"
	"service-courier/internal/router"
	"service-courier/internal/server"
	"service-courier/internal/service/courierapp"
	"service-courier/internal/service/deliveryapp"
	"service-courier/internal/worker"
	"syscall"
	"time"
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
	router := router.SetupRoute(
		healthHTTP,
		courHTTP,
		delHTTP,
	)

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
