package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Gazmasater/api"
	"github.com/Gazmasater/docs"
	"github.com/Gazmasater/internal/db"
	"github.com/Gazmasater/kafka"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Загрузка переменных окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Получение переменных окружения для подключения к базе данных
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Инициализация логгера
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Подключение к базе данных
	database, err := db.Connect(host, port, user, password, dbname)
	if err != nil {
		sugar.Fatalf("Failed to connect to database: %v", err)
	}

	// Создание роутера
	r := api.NewRouter(database, sugar)

	// Обновление информации о Swagger
	docs.SwaggerInfo.Title = "API MICROSERV"
	docs.SwaggerInfo.Description = "Это пример API для отправки сообщений."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	// Запуск консьюмера в отдельной горутине
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		kafka.StartConsumer(database)
	}()

	// Запуск сервера
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	sugar.Infof("Starting server on port %s", port)
	go func() {
		if err := http.ListenAndServe(":"+port, r); err != nil {
			sugar.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидание завершения работы консьюмера
	wg.Wait()
}
