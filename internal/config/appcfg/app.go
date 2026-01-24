package appcfg

import (
	"os"
)

type AppEnv struct {
	AppHost   string
	AppPort   string
	OrderHost string
	OrderPort string
	TimeCheck string
	TimePoll  string
	Refill    string
	Limit     string
}

// SetupAppEnv - парсер env переменных
func SetupAppEnv() *AppEnv {
	host := os.Getenv("COURIER_LOCALHOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("COURIER_LOCALPORT")
	if port == "" {
		port = "8080"
	}
	return &AppEnv{
		AppHost:   host,
		AppPort:   port,
		OrderHost: os.Getenv("ORDER_HOST"),
		OrderPort: os.Getenv("ORDER_GRPC_PORT"),
		TimeCheck: os.Getenv("TIME_CHECK"),
		TimePoll:  os.Getenv("TIME_POLL"),
		Refill:    os.Getenv("REFILL"),
		Limit:     os.Getenv("LIMIT"),
	}
}
