FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .
RUN swag init

# Собираем приложение
RUN go build -o main .

# Используем более актуальный образ Ubuntu для выполнения приложения
FROM ubuntu:22.04

# Устанавливаем рабочий каталог
WORKDIR /app

# Копируем все файлы и директории из этапа сборки
COPY --from=builder /app /app

# Устанавливаем необходимые зависимости
RUN apt-get update && apt-get install -y libc6 wget ca-certificates

# Запускаем приложение
CMD ["./main"]
