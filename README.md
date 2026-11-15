# Кинотека

Каталог из 1000+ карточек фильмов и мультиков. Свободные ленты (Creative Commons и общественное достояние) воспроизводятся прямо на сайте с потоковой подкачкой.

## Файлы

- `index.html` — лендинг/главная страница.
- `filmotheka.html` — сам каталог с фильмами, мультиками, поиском, пагинацией и плеером.
- `Dockerfile`, `docker-compose.yml`, `nginx.conf` — развёртывание через Docker.
- `.dockerignore`, `.gitignore` — системные файлы.

## Запуск

Просто откройте `index.html` в браузере — сайт автономен и работает без сервера.

### Через Docker

```bash
docker compose up -d --build
```

Сайт будет доступен на http://localhost:8080/ (лендинг) и http://localhost:8080/filmotheka.html (каталог).

## О контенте

- 17+ настоящих лент в свободном доступе (мультфильмы Blender Foundation, классика в общественном достоянии) — они проигрываются с потоковой подкачкой из Internet Archive, у каждой есть кнопка «Скачать».
- Остальные карточки — каталог с описаниями, обложками и рейтингами.

## Технологии

- HTML, CSS, JavaScript (без зависимостей и сборки).
- Локальное хранилище для личных данных в браузере (localStorage).

---

## Бэкенд (папка `backend/`)

REST API каталога фильмов на Go + PostgreSQL + Swagger.

### Необходимые инструменты

- Docker
- Go 1.17+
- PostgreSQL

### Установка и запуск

Создайте файл `.env` в корне проекта `backend/` и укажите переменные среды:

```env
POSTGRES_USER=filmotheka_user
POSTGRES_PASSWORD=filmotheka_pass
POSTGRES_DB=filmotheka_db
POSTGRES_PORT=5432
POSTGRES_HOST=db
SERVER_PORT=8080
```

Запустите приложение:

```bash
docker compose up --build
```

API будет доступен на http://localhost:8080.

### Тестирование

```bash
go test ./... -v -coverprofile=coverage/cover.out && go tool cover -html=coverage/cover.out -o coverage/cover.html
```

Затем откройте `coverage/cover.html` в браузере.

### Swagger

Перед запуском перегенерируйте документацию, если меняли аннотации:

```bash
swag init -g cmd/filmotheka/main.go
```

### Документация

Документация Swagger доступна по адресу http://localhost:8080/swagger/index.html

### Эндпоинты API

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/health` | Проверка работоспособности |
| GET | `/api/films` | Список фильмов (`?genre=` и `?limit=`) |
| POST | `/api/films` | Создать фильм |
| GET | `/api/films/{id}` | Получить фильм по id |
| PUT | `/api/films/{id}` | Обновить фильм |
| DELETE | `/api/films/{id}` | Удалить фильм |