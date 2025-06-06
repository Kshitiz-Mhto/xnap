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

	FromEmail         string
	FromEmailPassword string
	FromEmailSMTP     string
	SMTPAddress       string
	OWNER_EMAIL       string

	XNAP_DB                  string
	XNAP_LOGS_TABLE          string
	XNAP_BACKUP_TABLE        string
	BACKUP_OR_RESTORE_STATUS string

	AES_KEY_SIZE string
	KEY_LOCATION string
	KEY_FILE     string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		MySQL_DB_HOST:            getEnv("MySQL_DB_HOST", "http://localhost"),
		MySQL_DB_PORT:            getEnv("MySQL_DB_PORT", "3360"),
		MySQL_DB_USER:            getEnv("MySQL_DB_USER", "root"),
		MySQL_DB_PASSWORD:        getEnv("MySQL_DB_PASSWORD", ""),
		POSTGRES_DB_USER:         getEnv("POSTGRES_DB_USER", "root"),
		POSTGRES_DB_PASSWORD:     getEnv("POSTGRES_DB_PASSWORD", ""),
		POSTGRES_DB_HOST:         getEnv("POSTGRES_DB_HOST", "http://localhost"),
		POSTGRES_DB_PORT:         getEnv("POSTGRES_DB_PORT", "5432"),
		FromEmail:                getEnv("FROM_EMAIL", ""),
		FromEmailPassword:        getEnv("FROM_EMAIL_PASSWORD", ""),
		FromEmailSMTP:            getEnv("FROM_EMAIL_SMTP", "smtp.gmail.com"),
		SMTPAddress:              getEnv("SMTP_ADDR", "smtp.gmail.com:587"),
		XNAP_DB:                  getEnv("LOG_DB", "xnap_db"),
		BACKUP_OR_RESTORE_STATUS: getEnv("BACKUP_OR_RESTORE_STATUS", ""),
		OWNER_EMAIL:              getEnv("OWNER_EMAIL", ""),
		XNAP_LOGS_TABLE:          getEnv("XNAP_LOGS_TABLE", "logs"),
		XNAP_BACKUP_TABLE:        getEnv("XNAP_BACKUP_TABLE", "backups"),
		AES_KEY_SIZE:             getEnv("AES_KEY_SIZE", "32"),
		KEY_LOCATION:             getEnv("KEY_LOCATION", ".xnap"),
		KEY_FILE:                 getEnv("KEY_FILE", "key"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
