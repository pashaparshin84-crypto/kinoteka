// Package config загружает настройки приложения из окружения.
package config

import (
	"os"
)

// Config хранит параметры подключения к БД и HTTP-сервера.
type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresPort     string
	PostgresHost     string
	ServerPort       string
}

// Load читает переменные окружения и возвращает конфигурацию.
// Все значения имеют разумные значения по умолчанию.
func Load() Config {
	return Config{
		PostgresUser:     getenv("POSTGRES_USER", "filmotheka_user"),
		PostgresPassword: getenv("POSTGRES_PASSWORD", "filmotheka_pass"),
		PostgresDB:       getenv("POSTGRES_DB", "filmotheka_db"),
		PostgresPort:     getenv("POSTGRES_PORT", "5432"),
		PostgresHost:     getenv("POSTGRES_HOST", "db"),
		ServerPort:       getenv("SERVER_PORT", "8080"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}