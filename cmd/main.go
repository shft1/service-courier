package main

import (
	"context"
	"service-courier/internal/bootstrap"
	"service-courier/internal/route"
)

func main() {
	env := bootstrap.SetupEnv()
	bootstrap.CliHandler(context.Background(), env)
	bootstrap.StartServerGraceful(route.SetupRoute(), env)
}
