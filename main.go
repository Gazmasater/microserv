package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Gazmasater/api"
	"github.com/Gazmasater/docs"
	"github.com/Gazmasater/internal/db"
	"github.com/Gazmasater/kafka"
	"go.uber.org/zap"
)

func main() {
	// Информация для Swagger документации
	docs.SwaggerInfo.Title = "API MICROSERV"
	docs.SwaggerInfo.Description = "Это пример API для отправки сообщений."
	docs.SwaggerInfo.Version = "1.1"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Параметры подключения к базе данных
	dbHost := "postgres"
	dbPort := "5432"
	dbUser := "postgres"
	dbPassword := "qwert"
	dbName := "microserv"
	port := "8080"

	// Инициализация логгера
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Создаем контекст
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Подключение к базе данных
	database, err := db.Connect(ctx, dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		sugar.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer database.Close()

	// Запуск миграций базы данных
	// if err := db.RunMigrations(database); err != nil {
	// 	sugar.Fatalf("Не удалось выполнить миграции: %v", err)
	// }

	// Инициализация роутера
	r := api.NewRouter(database, sugar)

	// Создаем канал для остановки Kafka consumer
	stopKafka := make(chan struct{})

	// Запуск Kafka consumer в отдельной горутине
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		kafka.StartConsumer(ctx, database, stopKafka) // Передаем контекст и сигнал остановки
	}()

	// Обработка корректного завершения работы
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		sugar.Infof("Запуск сервера на IP %s и порту %s", srv.Addr, port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("Не удалось запустить сервер: %v", err)
		}
	}()

	// Ожидание сигнала для корректного завершения
	<-stop
	sugar.Info("Остановка сервера...")

	// Завершаем работу сервера
	if err := srv.Shutdown(context.Background()); err != nil {
		sugar.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	// Отправляем сигнал остановки Kafka consumer
	close(stopKafka)

	// Отменяем контекст для Kafka consumer
	cancel()

	// Ждем завершения всех горутин
	wg.Wait()
	sugar.Info("Сервер корректно остановлен")
}
