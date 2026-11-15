package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHealth проверяет, что эндпоинт /health отвечает 200 со статусом ok.
func TestHealth(t *testing.T) {
	h := New(nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	h.Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ожидали 200, получили %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("некорректный JSON: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("ожидали status=ok, получили %q", body["status"])
	}
}

// TestPathID проверяет извлечение id из пути /api/films/{id}.
func TestPathID(t *testing.T) {
	cases := []struct {
		path string
		want int64
		ok   bool
	}{
		{"/api/films/42", 42, true},
		{"/api/films/7", 7, true},
		{"/api/films/abc", 0, false},
		{"/api/films/0", 0, false},
		{"/api/films/-3", 0, false},
		{"/api/films/", 0, false},
		{"/api/films", 0, false},
	}
	for _, c := range cases {
		req := httptest.NewRequest(http.MethodGet, c.path, nil)
		id, ok := pathID(req)
		if ok != c.ok || id != c.want {
			t.Errorf("pathID(%q) = (%d,%v), ожидали (%d,%v)", c.path, id, ok, c.want, c.ok)
		}
	}
}

// TestWriteError проверяет вывод ошибки в формате JSON.
func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()
	writeError(rec, http.StatusBadRequest, "плохой запрос")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("ожидали 400, получили %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Fatalf("ожидали application/json, получили %q", ct)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("некорректный JSON: %v", err)
	}
	if body["error"] != "плохой запрос" {
		t.Fatalf("ожидали error=плохой запрос, получили %q", body["error"])
	}
}