# ---- Стадия сборки ----
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Копируем модули и подтягиваем зависимости (кэширование слоёв)
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники и собираем бинарник
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /filmotheka ./cmd/filmotheka

# ---- Стадия запуска ----
FROM alpine:3.20

RUN apk add --no-cache ca-certificates && adduser -D -u 10001 appuser
USER appuser

WORKDIR /app
COPY --from=builder /filmotheka /app/filmotheka

EXPOSE 8080
CMD ["/app/filmotheka"]