package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"service-courier/internal/cli"
	"service-courier/internal/config"
	"service-courier/internal/db"
	"service-courier/internal/handler"
	"service-courier/internal/repository"
	"service-courier/internal/router"
	"service-courier/internal/server"
	"service-courier/internal/service"
	"syscall"
)

func main() {
	sysCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	env := config.SetupEnv()
	pool := db.InitPool(sysCtx, env)
	defer pool.Close()

	healthHandler := handler.NewHealthHandler()

	courierRepository := repository.NewCourierRepository(pool)
	courierService := service.NewCourierService(courierRepository)
	courierHandler := handler.NewCourierHandler(courierService)

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
