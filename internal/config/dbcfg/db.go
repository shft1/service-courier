package dbcfg

import (
	"os"
)

type DataBaseEnv struct {
	DBUser string
	DBPass string
	DBName string
	DBHost string
	DBPort string
}

// SetupDataBaseEnv - парсер env переменных
func SetupDataBaseEnv() *DataBaseEnv {
	return &DataBaseEnv{
		DBUser: os.Getenv("POSTGRES_USER"),
		DBPass: os.Getenv("POSTGRES_PASSWORD"),
		DBName: os.Getenv("POSTGRES_DB"),
		DBHost: os.Getenv("POSTGRES_HOST"),
		DBPort: os.Getenv("POSTGRES_LOCALPORT"),
	}
}
