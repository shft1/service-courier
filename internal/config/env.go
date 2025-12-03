package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Port      string
	DBUser    string
	DBPass    string
	DBName    string
	DBPort    string
	DBHost    string
	TimeCheck string
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
		Port:      port,
		DBUser:    os.Getenv("POSTGRES_USER"),
		DBPass:    os.Getenv("POSTGRES_PASSWORD"),
		DBName:    os.Getenv("POSTGRES_DB"),
		DBPort:    os.Getenv("POSTGRES_PORT"),
		DBHost:    os.Getenv("POSTGRES_HOST"),
		TimeCheck: os.Getenv("TIME_CHECK"),
	}
}
