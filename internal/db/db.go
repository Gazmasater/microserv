package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func Connect(ctx context.Context, host, port, user, password, dbname string) (*sql.DB, error) {
	// Формирование строки подключения
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	fmt.Println("Connecting to", host)
	fmt.Printf("Connecting to database with the following details:\nHost: %s\nPort: %s\nUser: %s\nPassword: %s\nDBName: %s\n",
		host, port, user, password, dbname)

	// Подготовка подключения к базе данных
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Установка таймаута для проверки соединения
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Проверка соединения
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")

	return db, nil
}
