package kafka

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/Gazmasater/internal/models"
	"github.com/IBM/sarama"
)

func StartConsumer(db *sql.DB) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, config)
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

	var wg sync.WaitGroup

	// Обработка сообщений
	fmt.Println("Начало обработки сообщений...")
	for msg := range partitionConsumer.Messages() {
		wg.Add(1)
		go func(msg *sarama.ConsumerMessage) {
			defer wg.Done()

			fmt.Printf("Получено сообщение: %s\n", string(msg.Value))

			var message models.Message
			if err := json.Unmarshal(msg.Value, &message); err != nil {
				log.Printf("Не удалось распарсить сообщение: %v", err)
				return
			}

			// Валидация сообщения
			if err := models.ValidateMessage(&message); err != nil {
				log.Printf("Ошибка валидации: %v", err)
				return
			}

			// Пометка сообщения как обработанного
			message.Status = models.StatusProcessed

			// Сохранение сообщения в базе данных
			if err := models.SaveMessage(db, &message); err != nil {
				log.Printf("Не удалось сохранить сообщение в базу данных: %v", err)
			}
		}(msg)
	}

	// Ожидание завершения всех goroutine
	wg.Wait()
	fmt.Println("Завершение обработки сообщений.")
}
