// Package handlers реализует HTTP-обработчики REST API каталога фильмов.
package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"filmotheka/internal/models"
)

// Handlers группирует зависимости REST-обработчиков.
type Handlers struct {
	DB *sql.DB
}

// New возвращает новый экземпляр обработчиков.
func New(db *sql.DB) *Handlers {
	return &Handlers{DB: db}
}

// Health отвечает на запросы проверки работоспособности сервиса.
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ListFilms возвращает список фильмов. Поддерживает query-параметры
// genre и limit для фильтрации.
// @Summary Список фильмов
// @Description Возвращает список фильмов каталога с опциональной фильтрацией по жанру.
// @Tags films
// @Produce json
// @Param genre query string false "Фильтр по жанру"
// @Param limit query int false "Ограничение количества (по умолчанию 100)"
// @Success 200 {array} models.Film
// @Router /api/films [get]
func (h *Handlers) ListFilms(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get("genre")
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "100"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 100
	}
	if limitInt > 500 {
		limitInt = 500
	}

	query := `SELECT id, title, genre, year, rating, description, video_url, is_cartoon, created_at
	          FROM films`
	args := []interface{}{}
	if genre != "" {
		query += ` WHERE genre ILIKE $1`
		args = append(args, "%"+genre+"%")
	}
	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(len(args)+1)
	args = append(args, limitInt)

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		serverError(w, "не удалось получить список фильмов", err)
		return
	}
	defer rows.Close()

	films := []models.Film{}
	for rows.Next() {
		var f models.Film
		if err := rows.Scan(&f.ID, &f.Title, &f.Genre, &f.Year, &f.Rating,
			&f.Description, &f.VideoURL, &f.IsCartoon, &f.CreatedAt); err != nil {
			serverError(w, "ошибка чтения строки", err)
			return
		}
		films = append(films, f)
	}
	writeJSON(w, http.StatusOK, films)
}

// CreateFilm добавляет новый фильм в каталог.
// @Summary Добавить фильм
// @Description Создаёт новый фильм в каталоге и возвращает его с присвоенным id.
// @Tags films
// @Accept json
// @Produce json
// @Param film body models.Film true "Данные нового фильма"
// @Success 201 {object} models.Film
// @Failure 400 {object} map[string]string
// @Router /api/films [post]
func (h *Handlers) CreateFilm(w http.ResponseWriter, r *http.Request) {
	var f models.Film
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}
	f.Title = strings.TrimSpace(f.Title)
	if f.Title == "" {
		writeError(w, http.StatusBadRequest, "название фильма обязательно")
		return
	}

	const q = `INSERT INTO films (title, genre, year, rating, description, video_url, is_cartoon)
	           VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, created_at`
	err := h.DB.QueryRow(q, f.Title, f.Genre, f.Year, f.Rating,
		f.Description, f.VideoURL, f.IsCartoon).Scan(&f.ID, &f.CreatedAt)
	if err != nil {
		serverError(w, "не удалось сохранить фильм", err)
		return
	}
	writeJSON(w, http.StatusCreated, f)
}

// GetFilm возвращает один фильм по id.
// @Summary Получить фильм
// @Description Возвращает фильм по его идентификатору.
// @Tags films
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {object} models.Film
// @Failure 404 {object} map[string]string
// @Router /api/films/{id} [get]
func (h *Handlers) GetFilm(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "некорректный id")
		return
	}
	var f models.Film
	const q = `SELECT id, title, genre, year, rating, description, video_url, is_cartoon, created_at
	           FROM films WHERE id=$1`
	err := h.DB.QueryRow(q, id).Scan(&f.ID, &f.Title, &f.Genre, &f.Year, &f.Rating,
		&f.Description, &f.VideoURL, &f.IsCartoon, &f.CreatedAt)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "фильм не найден")
		return
	}
	if err != nil {
		serverError(w, "ошибка получения фильма", err)
		return
	}
	writeJSON(w, http.StatusOK, f)
}

// UpdateFilm обновляет существующий фильм.
// @Summary Обновить фильм
// @Description Обновляет данные фильма по его id.
// @Tags films
// @Accept json
// @Param id path int true "ID фильма"
// @Param film body models.Film true "Новые данные"
// @Success 200 {object} models.Film
// @Failure 404 {object} map[string]string
// @Router /api/films/{id} [put]
func (h *Handlers) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "некорректный id")
		return
	}
	var f models.Film
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}
	const q = `UPDATE films SET title=$2, genre=$3, year=$4, rating=$5, description=$6, video_url=$7, is_cartoon=$8
	           WHERE id=$1 RETURNING id, created_at`
	err := h.DB.QueryRow(q, id, f.Title, f.Genre, f.Year, f.Rating,
		f.Description, f.VideoURL, f.IsCartoon).Scan(&f.ID, &f.CreatedAt)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "фильм не найден")
		return
	}
	if err != nil {
		serverError(w, "не удалось обновить фильм", err)
		return
	}
	writeJSON(w, http.StatusOK, f)
}

// DeleteFilm удаляет фильм по id.
// @Summary Удалить фильм
// @Description Удаляет фильм из каталога по его id.
// @Tags films
// @Param id path int true "ID фильма"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /api/films/{id} [delete]
func (h *Handlers) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "некорректный id")
		return
	}
	res, err := h.DB.Exec(`DELETE FROM films WHERE id=$1`, id)
	if err != nil {
		serverError(w, "не удалось удалить фильм", err)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "фильм не найден")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// pathID извлекает числовой id из пути /api/films/{id}.
func pathID(r *http.Request) (int64, bool) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 0 {
		return 0, false
	}
	raw := parts[len(parts)-1]
	if raw == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func serverError(w http.ResponseWriter, msg string, err error) {
	log.Printf("%s: %v", msg, err)
	writeError(w, http.StatusInternalServerError, msg)
}