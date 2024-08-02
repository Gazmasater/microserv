package kafka

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Gazmasater/internal/models"
	"github.com/Gazmasater/pkg/logger"
	"github.com/IBM/sarama"
)

func StartConsumer(ctx context.Context, db *sql.DB, stop <-chan struct{}) {
	brokerList := []string{"kafka:9092"}

	// Получаем экземпляр логгера
	sugar := logger.GetLogger()

	// Создание контекста с таймаутом для подключения к Kafka

	connectCtx, cancel := context.WithTimeout(ctx, models.ConnectionTimeoutDuration)
	defer cancel()

	if err := checkKafkaConnection(connectCtx, brokerList); err != nil {
		sugar.Fatalf("Не удалось подключиться к Kafka: %v", err)
	} else {
		sugar.Info("Успешно подключено к Kafka")
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	sugar.Infof("StartConsumer %v", brokerList)

	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		sugar.Fatalf("Не удалось запустить Kafka consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			sugar.Errorf("Ошибка при закрытии консюмера: %v", err)
		}
	}()

	sugar.Info("Запуск потребления партиции...")
	partitionConsumer, err := consumer.ConsumePartition("dbserver1.public.msg", 0, sarama.OffsetNewest)
	if err != nil {
		sugar.Fatalf("Не удалось запустить partition consumer: %v", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			sugar.Errorf("Ошибка при закрытии partition consumer: %v", err)
		}
	}()

	sugar.Info("Начало обработки сообщений...")
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			sugar.Infof("Получено сообщение: %s", string(msg.Value))

			var message models.Message
			if err := json.Unmarshal(msg.Value, &message); err != nil {
				sugar.Errorf("Не удалось распарсить сообщение: %v", err)
				continue
			}

			if message.Status_2 == models.StatusProcessed {
				sugar.Info("Сообщение со статусом 'processed', пропуск...")
				continue
			}

			message.Status_2 = models.StatusProcessed
			sugar.Infof("KAFKA message %+v", message)
			if err := models.SaveMessage2(db, &message); err != nil {
				sugar.Errorf("Не удалось сохранить сообщение в базу данных2: %v", err)
			} else {
				sugar.Info("Сообщение успешно сохранено в базу данных.")
			}

		case <-stop:
			sugar.Info("Kafka consumer завершает работу...")
			return // Выходим из функции при получении сигнала остановки
		case <-ctx.Done():
			sugar.Info("Контекст завершён, Kafka consumer завершает работу...")
			return // Выходим из функции при завершении контекста
		}
	}
}

func checkKafkaConnection(ctx context.Context, brokerList []string) error {
	config := sarama.NewConfig()
	var client sarama.Client
	var err error

	sugar := logger.GetLogger() // Получаем экземпляр логгера

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
			sugar.Errorf("Не удалось подключиться к Kafka: %v", err)
			time.Sleep(2 * time.Second) // ожидание 2 секунды между попытками
		}
	}
}
