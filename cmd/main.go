package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"service-courier/internal/cli"
	"service-courier/internal/config"
	"service-courier/internal/db/postgre"
	courierHandler "service-courier/internal/handler/courier"
	deliveryHandler "service-courier/internal/handler/delivery"
	healthHandler "service-courier/internal/handler/health"
	courierRepository "service-courier/internal/repository/courier"
	deliveryRepository "service-courier/internal/repository/delivery"
	"service-courier/internal/router"
	"service-courier/internal/server"
	courierService "service-courier/internal/service/courier"
	deliveryService "service-courier/internal/service/delivery"
	"service-courier/internal/worker"
	"syscall"
	"time"
)

func main() {
	sysCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	env := config.SetupEnv()
	pool := postgre.InitPool(sysCtx, env)
	defer pool.Close()
	txManager := postgre.NewTxManagerPostgre(pool)

	healthHandler := healthHandler.NewHealthHandler()

	courierRepository := courierRepository.NewCourierRepository(pool, txManager)
	courierService := courierService.NewCourierService(courierRepository)
	courierHandler := courierHandler.NewCourierHandler(courierService)

	deliveryRepository := deliveryRepository.NewDeliveryRepository(pool, txManager)
	deliveryService := deliveryService.NewDeliveryService(deliveryRepository, courierRepository, txManager)
	deliveryHandler := deliveryHandler.NewDeliveryHandler(deliveryService)

	router := router.SetupRoute(
		healthHandler,
		courierHandler,
		deliveryHandler,
	)

	checkPeriod := time.Second * 10
	deliveryChecker := worker.NewDeliveryMonitor(checkPeriod, deliveryService)
	go deliveryChecker.Start(sysCtx)

	cmd := cli.CliHandler(env)
	if err := cmd.Run(sysCtx, os.Args); err != nil {
		fmt.Println(err)
	}

	server.StartServerGraceful(sysCtx, router, pool, env)
}
