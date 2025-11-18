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
	healthHandler "service-courier/internal/handler/health"
	courierRepo "service-courier/internal/repository/courier"
	"service-courier/internal/router"
	"service-courier/internal/server"
	courierService "service-courier/internal/service/courier"
	"syscall"
)

func main() {
	sysCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	env := config.SetupEnv()
	pool := postgre.InitPool(sysCtx, env)
	defer pool.Close()

	healthHandler := healthHandler.NewHealthHandler()

	courierRepository := courierRepo.NewCourierRepository(pool)
	courierService := courierService.NewCourierService(courierRepository)
	courierHandler := courierHandler.NewCourierHandler(courierService)

	router := router.SetupRoute(
		healthHandler,
		courierHandler,
	)

	cmd := cli.CliHandler(env)
	if err := cmd.Run(sysCtx, os.Args); err != nil {
		fmt.Println(err)
	}

	server.StartServerGraceful(sysCtx, router, pool, env)
}
