package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Gazmasater/api"
	"github.com/Gazmasater/docs"
	"github.com/Gazmasater/internal/db"
	"github.com/Shopify/sarama"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Загрузка переменных окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Получение переменных окружения для подключения к базе данных и Kafka
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	kafkaBrokers := []string{os.Getenv("KAFKA_BROKERS")}
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	// Инициализация логгера
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Подключение к базе данных
	database, err := db.Connect(host, port, user, password, dbname)
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

	// Конфигурация и создание Kafka-клиента
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Version = sarama.V2_8_0_0 // Укажите нужную версию Kafka
	kafkaConfig.Consumer.Return.Errors = true

	kafkaClient, err := sarama.NewClient(kafkaBrokers, kafkaConfig)
	if err != nil {
		sugar.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer func() {
		if err := kafkaClient.Close(); err != nil {
			sugar.Fatalf("Failed to close Kafka client: %v", err)
		}
	}()

	// Запуск консьюмера в отдельной горутине
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		startConsumer(kafkaClient, kafkaTopic, database, sugar)
	}()

	// Запуск сервера
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	sugar.Infof("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		sugar.Fatalf("Failed to start server: %v", err)
	}

	// Ожидание завершения работы консьюмера
	wg.Wait()
}

func startConsumer(client sarama.Client, topic string, db *sql.DB, logger *zap.SugaredLogger) {
	// Создание конфигурации консьюмера
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0 // Укажите нужную версию Kafka
	config.Consumer.Return.Errors = true

	// Создание нового консьюмера
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		logger.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			logger.Fatalf("Failed to close Kafka consumer: %v", err)
		}
	}()

	// Получение партиций для топика
	partitions, err := client.Partitions(topic)
	if err != nil {
		logger.Fatalf("Failed to get partitions for topic %s: %v", topic, err)
	}

	// Создание горутин для потребления сообщений из всех партиций
	var wg sync.WaitGroup
	for _, partition := range partitions {
		wg.Add(1)
		go func(partition int32) {
			defer wg.Done()
			pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				logger.Fatalf("Failed to start consumer for partition %d: %v", partition, err)
			}
			defer pc.Close()

			for msg := range pc.Messages() {
				// Обработка сообщений
				fmt.Printf("Received message: %s\n", string(msg.Value))
				// Здесь можно добавить код для сохранения сообщений в базу данных
			}
		}(partition)
	}

	wg.Wait()
}
