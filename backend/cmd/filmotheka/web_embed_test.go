package main

import (
	"io/fs"
	"testing"
)

// TestWebFS проверяет, что веб-интерфейс встроен и index.html существует.
func TestWebFS(t *testing.T) {
	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		t.Fatalf("не удалось получить подкаталог web: %v", err)
	}
	if _, err := fs.Stat(sub, "index.html"); err != nil {
		t.Fatalf("index.html не найден во встроенном FS: %v", err)
	}
	t.Log("web/index.html встроен успешно")
}