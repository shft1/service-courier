package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Port string
}

func SetupEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		os.Setenv("PORT", port)
	}
	return &Env{
		Port: port,
	}
}
