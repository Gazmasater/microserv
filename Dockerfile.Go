# Используем официальный образ Go в качестве базового
FROM golang:1.22 AS builder

# Устанавливаем рабочий каталог
WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Устанавливаем swag CLI и goose
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Копируем весь исходный код и генерируем документацию с помощью swag
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
#COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Устанавливаем необходимые зависимости
RUN apt-get update && apt-get install -y libc6 wget ca-certificates

# Установка переменных окружения для goose
#ENV GOOSE_DRIVER=postgres
#ENV GOOSE_DBSTRING="user=postgres password=qwert dbname=microserv host=postgres port=5432 sslmode=disable"

# Проверка содержимого директории миграций и других файлов
#RUN ls -la /app/internal/db/migrations
RUN ls -la /app/docs
#RUN ls -la /usr/local/bin/goose

# Запускаем приложение
CMD ["./main"]
