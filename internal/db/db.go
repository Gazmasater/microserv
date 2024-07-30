package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Gazmasater/pkg/logger" // Импортируем ваш логгер

	_ "github.com/lib/pq"
)

func Connect(ctx context.Context, host, port, user, password, dbname string) (*sql.DB, error) {
	// Получаем экземпляр логгера
	sugar := logger.GetLogger()

	// Формирование строки подключения
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sugar.Infof("Connecting to database with the following details:\nHost: %s\nPort: %s\nUser: %s\nPassword: %s\nDBName: %s\n",
		host, port, user, password, dbname)

	// Подготовка подключения к базе данных
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Установка таймаута для проверки соединения
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Проверка соединения
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	sugar.Info("Successfully connected to database")

	return db, nil
}
