// Package db отвечает за подключение к PostgreSQL и миграции схемы.
package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // регистрация драйвера PostgreSQL

	"filmotheka/internal/config"
)

// Connect устанавливает соединение с PostgreSQL и применяет миграцию.
// Выполняет несколько попыток подключения, чтобы дождаться готовности БД.
func Connect(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB,
	)

	var db *sql.DB
	var err error

	for attempt := 1; attempt <= 15; attempt++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			if pingErr := db.Ping(); pingErr == nil {
				break
			}
		}
		log.Printf("БД ещё не готова (попытка %d/15), повтор через 2с...", attempt)
		if db != nil {
			db.Close()
		}
		time.Sleep(2 * time.Second)
		err = fmt.Errorf("connection failed")
	}

	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("ошибка миграции: %w", err)
	}

	log.Println("Подключение к PostgreSQL установлено")
	return db, nil
}

// migrate создаёт таблицу films, если её ещё нет.
func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS films (
			id          BIGSERIAL PRIMARY KEY,
			title       TEXT NOT NULL,
			genre       TEXT NOT NULL DEFAULT '',
			year        INT  NOT NULL DEFAULT 0,
			rating      DOUBLE PRECISION NOT NULL DEFAULT 0,
			description TEXT NOT NULL DEFAULT '',
			video_url   TEXT NOT NULL DEFAULT '',
			is_cartoon  BOOLEAN NOT NULL DEFAULT FALSE,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	return err
}