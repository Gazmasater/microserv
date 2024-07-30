# Используем официальный образ Go в качестве базового
FROM golang:1.22 AS builder

# Устанавливаем рабочий каталог
WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Устанавливаем swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

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

# Устанавливаем необходимые зависимости
RUN apt-get update && apt-get install -y libc6 wget ca-certificates

# Запускаем приложение
CMD ["./main"]
