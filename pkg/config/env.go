package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MySQL_DB_HOST     string
	MySQL_DB_PORT     string
	MySQL_DB_USER     string
	MySQL_DB_PASSWORD string

	POSTGRES_DB_USER     string
	POSTGRES_DB_PASSWORD string
	POSTGRES_DB_HOST     string
	POSTGRES_DB_PORT     string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		MySQL_DB_HOST:        getEnv("MySQL_DB_HOST", "http://localhost"),
		MySQL_DB_PORT:        getEnv("MySQL_DB_PORT", "3360"),
		MySQL_DB_USER:        getEnv("MySQL_DB_USER", "root"),
		MySQL_DB_PASSWORD:    getEnv("MySQL_DB_PASSWORD", ""),
		POSTGRES_DB_USER:     getEnv("POSTGRES_DB_USER", "root"),
		POSTGRES_DB_PASSWORD: getEnv("POSTGRES_DB_PASSWORD", ""),
		POSTGRES_DB_HOST:     getEnv("POSTGRES_DB_HOST", "http://localhost"),
		POSTGRES_DB_PORT:     getEnv("POSTGRES_DB_PORT", "5432"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
