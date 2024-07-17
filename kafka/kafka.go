package kafka

import (
	"database/sql"
	"encoding/json"
	"log"
	"sync"

	"github.com/Gazmasater/internal/models"
	"github.com/IBM/sarama"
)

// ProduceMessage отправляет сообщение в Kafka
func ProduceMessage(message interface{}) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Printf("Failed to start Sarama producer: %v", err)
		return err
	}
	defer producer.Close()

	msg, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return err
	}

	kafkaMessage := &sarama.ProducerMessage{
		Topic: "your_topic_name", // Замените на ваше имя топика
		Value: sarama.StringEncoder(msg),
	}

	_, _, err = producer.SendMessage(kafkaMessage)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	return nil
}

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
