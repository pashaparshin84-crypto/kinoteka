// cmd/filmotheka — точка входа REST API каталога фильмов.
//
// Подключается к PostgreSQL, настраивает маршруты и запускает HTTP-сервер.
package main

import (
	"log"
	"net/http"

	"github.com/swaggo/http-swagger/v2"

	_ "filmotheka/docs" // сгенерированная документация Swagger
	"filmotheka/internal/config"
	"filmotheka/internal/db"
	"filmotheka/internal/handlers"
)

// @title Кинотека API
// @version 1.0
// @description REST API каталога фильмов и мультиков на Go + PostgreSQL.

// @contact.name Кинотека
// @contact.url http://localhost:8080

// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.Load()

	conn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("FATAL: %v", err)
	}
	defer conn.Close()

	h := handlers.New(conn)
	mux := http.NewServeMux()

	// Swagger UI: http://localhost:8080/swagger/index.html
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Маршруты
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /api/films", h.ListFilms)
	mux.HandleFunc("POST /api/films", h.CreateFilm)
	mux.HandleFunc("GET /api/films/{id}", h.GetFilm)
	mux.HandleFunc("PUT /api/films/{id}", h.UpdateFilm)
	mux.HandleFunc("DELETE /api/films/{id}", h.DeleteFilm)

	addr := ":" + cfg.ServerPort
	log.Printf("Кинотека API запущен на %s", addr)
	if err := http.ListenAndServe(addr, logMiddleware(mux)); err != nil {
		log.Fatalf("FATAL: сервер завершился: %v", err)
	}
}

// logMiddleware логирует каждый входящий HTTP-запрос.
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}