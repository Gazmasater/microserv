package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Gazmasater/pkg/logger"

	"github.com/Gazmasater/api"
	"github.com/Gazmasater/docs"
	"github.com/Gazmasater/internal/db"
	"github.com/Gazmasater/kafka"
)

func main() {

	logger.Init()
	sugar := logger.GetLogger()

	// Информация для Swagger документации
	docs.SwaggerInfo.Title = "API MICROSERV"
	docs.SwaggerInfo.Description = "Это пример API для отправки сообщений."
	docs.SwaggerInfo.Version = "1.1"
	docs.SwaggerInfo.Host = os.Getenv("HOST") + ":" + os.Getenv("PORT")
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Параметры подключения к базе данных

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("PORT")

	// Создаем контекст
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Подключение к базе данных
	database, err := db.Connect(ctx, dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		sugar.Fatalf("Не удалось подключиться к базе данных: %v", err)

	}
	defer database.Close()

	// Инициализация роутера
	r := api.NewRouter(database, sugar)

	stopKafka := make(chan struct{})

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
