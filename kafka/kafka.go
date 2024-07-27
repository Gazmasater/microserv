package kafka

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Gazmasater/internal/models"
	"github.com/IBM/sarama"
)

func StartConsumer(db *sql.DB, stop <-chan struct{}) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

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

			if message.Status == models.StatusProcessed {
				fmt.Println("Сообщение со статусом 'processed', пропуск...")
				continue
			}

			message.Status = models.StatusProcessed

			if err := models.SaveMessage(db, &message); err != nil {
				log.Printf("Не удалось сохранить сообщение в базу данных: %v", err)
			} else {
				fmt.Println("Сообщение успешно сохранено в базу данных.")
			}

		case <-stop:
			log.Println("Kafka consumer завершает работу...")
			return // Выходим из функции при получении сигнала остановки
		}
	}
}
