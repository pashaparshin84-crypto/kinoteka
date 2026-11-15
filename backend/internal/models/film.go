// Package models определяет структуры данных приложения.
package models

import "time"

// Film представляет запись о фильме в каталоге.
// @Description Фильм в каталоге Киноте ки
type Film struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Genre       string    `json:"genre"`
	Year        int       `json:"year"`
	Rating      float64   `json:"rating"`
	Description string    `json:"description"`
	VideoURL    string    `json:"video_url,omitempty"`
	IsCartoon   bool      `json:"is_cartoon"`
	CreatedAt   time.Time `json:"created_at"`
}