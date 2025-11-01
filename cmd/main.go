package main

import (
	"context"
	"net/http"
	"serviceDelivery/configs"
	"serviceDelivery/internal/route"
)

func main() {
	env := configs.SetupEnv()
	configs.CliHandler(context.Background(), env)
	http.ListenAndServe(":"+env.Port, route.SetupRoute())
}
