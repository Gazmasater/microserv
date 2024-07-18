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
	// Создание конфигурации клиента
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0 // Укажите нужную версию Kafka
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true // Включаем авто-коммит смещений
	config.ClientID = "your-consumer-group"

	// Создание нового клиента
	client, err := sarama.NewClient([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Не удалось создать клиента Kafka: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalf("Не удалось закрыть клиента Kafka: %v", err)
		}
	}()

	// Получение метаданных топиков
	metaFetcher := sarama.NewMetadataFetcher(client)
	meta, err := metaFetcher.FetchMetadata([]string{"my_topic"})
	if err != nil {
		log.Fatalf("Не удалось получить метаданные топика: %v", err)
	}

	var wg sync.WaitGroup

	// Создаем один консьюмер для всех партиций
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Fatalf("Не удалось создать консьюмера: %v", err)
	}
	defer consumer.Close()

	// Получаем партиции для топика "my_topic"
	for _, partition := range meta.Topics["my_topic"].Partitions {
		wg.Add(1)
		go func(partitionID int32) {
			defer wg.Done()

			partitionConsumer, err := consumer.ConsumePartition("my_topic", partitionID, sarama.OffsetNewest)
			if err != nil {
				log.Printf("Не удалось запустить partition consumer для партиции %d: %v", partitionID, err)
				return
			}
			defer partitionConsumer.Close()

			for msg := range partitionConsumer.Messages() {
				// Обработка сообщения
				var message models.Message
				if err := json.Unmarshal(msg.Value, &message); err != nil {
					log.Printf("Не удалось распарсить сообщение: %v", err)
					continue
				}

				// Валидация сообщения
				if err := models.ValidateMessage(&message); err != nil {
					log.Printf("Ошибка валидации: %v", err)
					continue
				}

				message.Status = models.StatusProcessed

				// Сохранение сообщения в базе данных
				if err := models.SaveMessage(db, &message); err != nil {
					log.Printf("Не удалось сохранить сообщение в базу данных: %v", err)
				}
			}
		}(partition.ID)
	}

	// Ожидание завершения всех горутин
	wg.Wait()
}
