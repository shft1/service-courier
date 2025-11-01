package main

import (
	"context"
	"serviceDelivery/internal/bootstrap"
	"serviceDelivery/internal/route"
)

func main() {
	env := bootstrap.SetupEnv()
	bootstrap.CliHandler(context.Background(), env)
	bootstrap.StartServerGraceful(route.SetupRoute(), env)
}
