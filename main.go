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
	dbHost := "postgres"
	dbPort := "5432"
	dbUser := "postgres"
	dbPassword := "qwert"
	dbName := "microserv"
	port := "8080"

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	database, err := db.Connect(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		sugar.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.RunMigrations(database); err != nil {
		sugar.Fatalf("Failed to run migrations: %v", err)
	}

	r := api.NewRouter(database, sugar)

	docs.SwaggerInfo.Title = "API MICROSERV"
	docs.SwaggerInfo.Description = "Это пример API для отправки сообщений."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		kafka.StartConsumer(database) // Передаем только базу данных
	}()

	sugar.Infof("Starting server on port %s", port)
	go func() {
		if err := http.ListenAndServe(":"+port, r); err != nil {
			sugar.Fatalf("Failed to start server: %v", err)
		}
	}()

	wg.Wait()
}
