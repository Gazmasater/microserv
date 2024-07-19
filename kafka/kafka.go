package kafka

import (
	"database/sql"
	"encoding/json"
	"log"
	"sync"

	"github.com/Gazmasater/internal/models"
	"github.com/IBM/sarama"
)

func StartConsumer(db *sql.DB) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Не удалось запустить Kafka consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("my_topic", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Не удалось запустить partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	var wg sync.WaitGroup

	for msg := range partitionConsumer.Messages() {
		wg.Add(1)
		go func(msg *sarama.ConsumerMessage) {
			defer wg.Done()

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
		}(msg) // Передача msg напрямую, без &msg
	}
	wg.Wait()
}
