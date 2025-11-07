package main

import (
	"context"
	"service-courier/internal/bootstrap"
	"service-courier/internal/router"
)

func main() {
	env := bootstrap.SetupEnv()
	pool := bootstrap.InitPool(context.Background(), env)
	router := router.SetupRoute(pool)
	bootstrap.CliHandler(context.Background(), env)
	bootstrap.StartServerGraceful(router, pool, env)
}
