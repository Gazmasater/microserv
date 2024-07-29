package kafka

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Gazmasater/internal/models"
	"github.com/IBM/sarama"
)

func StartConsumer(ctx context.Context, db *sql.DB, stop <-chan struct{}) {
	brokerList := []string{"kafka:9092"}

	// Создание контекста с таймаутом для подключения к Kafka
	connectCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	if err := checkKafkaConnection(connectCtx, brokerList); err != nil {
		log.Fatalf("Не удалось подключиться к Kafka: %v", err)
	} else {
		log.Println("Успешно подключено к Kafka")
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	fmt.Println("StartConsumer", brokerList)

	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Fatalf("Не удалось запустить Kafka consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("Ошибка при закрытии консюмера: %v", err)
		}
	}()

	fmt.Println("Запуск потребления партиции...")
	partitionConsumer, err := consumer.ConsumePartition("dbserver1.public.msg", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Не удалось запустить partition consumer: %v", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Printf("Ошибка при закрытии partition consumer: %v", err)
		}
	}()

	fmt.Println("Начало обработки сообщений...")
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Получено сообщение: %s\n", string(msg.Value))

			var message models.Message
			if err := json.Unmarshal(msg.Value, &message); err != nil {
				log.Printf("Не удалось распарсить сообщение: %v", err)
				continue
			}

			if message.Status_2 == models.StatusProcessed {
				fmt.Println("Сообщение со статусом 'processed', пропуск...")
				continue
			}

			message.Status_2 = models.StatusProcessed
			fmt.Printf("KAFKA message %+v\n", message)
			if err := models.SaveMessage2(db, &message); err != nil {
				log.Printf("Не удалось сохранить сообщение в базу данных2: %v", err)
			} else {
				fmt.Println("Сообщение успешно сохранено в базу данных.")
			}

		case <-stop:
			log.Println("Kafka consumer завершает работу...")
			return // Выходим из функции при получении сигнала остановки
		case <-ctx.Done():
			log.Println("Контекст завершён, Kafka consumer завершает работу...")
			return // Выходим из функции при завершении контекста
		}
	}
}

func checkKafkaConnection(ctx context.Context, brokerList []string) error {
	config := sarama.NewConfig()
	var client sarama.Client
	var err error

	// Попытка подключения в течение времени, заданного в контексте
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			client, err = sarama.NewClient(brokerList, config)
			if err == nil {
				client.Close()
				return nil
			}
			fmt.Printf("Не удалось подключиться к Kafka: %v\n", err)
			time.Sleep(2 * time.Second) // ожидание 2 секунды между попытками
		}
	}
}
