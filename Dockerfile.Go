# Используем официальный образ Go в качестве базового
FROM golang:1.22 AS builder

# Устанавливаем рабочий каталог
WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код и собираем приложение
COPY . .
RUN go build -o main .

# Используем более актуальный образ Ubuntu для выполнения приложения
FROM ubuntu:22.04

# Устанавливаем рабочий каталог
WORKDIR /app

# Копируем скомпилированное приложение и миграции из предыдущего этапа
COPY --from=builder /app/main .
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations

# Устанавливаем необходимые зависимости
RUN apt-get update && apt-get install -y libc6

# Запускаем приложение
CMD ["./main"]
