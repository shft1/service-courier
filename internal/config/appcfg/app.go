package appcfg

import (
	"os"
)

type AppEnv struct {
	AppPort   string
	OrderHost string
	OrderPort string
	TimeCheck string
	TimePoll  string
}

// SetupAppEnv - парсер env переменных
func SetupAppEnv() *AppEnv {
	port := os.Getenv("COURIER_LOCALPORT")
	if port == "" {
		port = "8080"
		os.Setenv("COURIER_LOCALPORT", port)
	}
	return &AppEnv{
		AppPort:   port,
		OrderHost: os.Getenv("ORDER_HOST"),
		OrderPort: os.Getenv("ORDER_GRPC_PORT"),
		TimeCheck: os.Getenv("TIME_CHECK"),
		TimePoll:  os.Getenv("TIME_POLL"),
	}
}
