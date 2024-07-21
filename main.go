package main

import (
	"net/http"
	"sync"

	"github.com/Gazmasater/api"
	"github.com/Gazmasater/docs"
	"github.com/Gazmasater/internal/db"
	"github.com/Gazmasater/kafka"
	"go.uber.org/zap"
)

func main() {
	// Установка значений переменных окружения прямо в коде
	dbHost := "postgres"  // Укажите хост вашей базы данных
	dbPort := "5432"      // Укажите порт вашей базы данных
	dbUser := "postgres"  // Укажите пользователя вашей базы данных
	dbPassword := "qwert" // Укажите пароль вашей базы данных
	dbName := "microserv" // Укажите имя вашей базы данных
	//	kafkaBroker := "localhost:9092" // Укажите адрес Kafka брокера
	port := "8080" // Порт, на котором будет запущен ваш сервер

	// Инициализация логгера
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Подключение к базе данных
	database, err := db.Connect(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		sugar.Fatalf("Failed to connect to database: %v", err)
	}

	// Запуск миграций
	if err := db.RunMigrations(database); err != nil {
		sugar.Fatalf("Failed to run migrations: %v", err)
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
	sugar.Infof("Starting server on port %s", port)
	go func() {
		if err := http.ListenAndServe(":"+port, r); err != nil {
			sugar.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидание завершения работы консьюмера
	wg.Wait()
}
