package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Gazmasater/pkg/logger" // Импортируем ваш логгер

	_ "github.com/lib/pq"
)

// Connect устанавливает соединение с базой данных Postgres с логикой повторных попыток подключения
func Connect(ctx context.Context, host, port, user, password, dbname string) (*sql.DB, error) {
	// Получаем экземпляр логгера
	sugar := logger.GetLogger()

	// Формирование строки подключения
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sugar.Infof("Connecting to database with the following details:\nHost: %s\nPort: %s\nUser: %s\nPassword: %s\nDBName: %s\n",
		host, port, user, password, dbname)

	// Логика повторных попыток подключения
	retryInterval := 5 * time.Second
	var db *sql.DB
	var err error

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while trying to connect to database: %w", ctx.Err())
		default:
			// Подготовка подключения к базе данных
			db, err = sql.Open("postgres", psqlInfo)
			if err != nil {
				sugar.Errorf("Failed to open database connection: %v", err)
				time.Sleep(retryInterval)
				continue
			}

			// Установка таймаута для проверки соединения
			connCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			// Проверка соединения
			err = db.PingContext(connCtx)
			if err == nil {
				sugar.Info("Successfully connected to database")
				return db, nil
			}

			sugar.Errorf("Failed to ping database: %v", err)
			time.Sleep(retryInterval)
		}
	}
}
