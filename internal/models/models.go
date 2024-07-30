package models

import (
	"database/sql"

	"fmt"
	"log"
)

const (
	StatusProcessed = "processed" // Сообщение было обработано
	StatusPending   = "pending"   // Сообщение ожидает обработки
	StatusFailed    = "failed"    // Обработка сообщения завершилась неудачей
)

type Message struct {
	ID          int64  `json:"id"`
	Text        string `json:"text"`
	Status_1    string `json:"status_1"`
	Status_2    string `json:"status_2"`
	CreatedAt_1 int64  `json:"created_at1"`
	CreatedAt_2 int64  `json:"created_at2"`
}

type Message_Request struct {
	Text string `json:"text"` // Набранный текст
}

type Stats struct {
	PendingMessages   int `json:"pending_messages"`   //количество необработанных сообщений
	ProcessedMessages int `json:"processed_messages"` //количество обработанных сообщений
	TotalMessages     int `json:"total_messages"`     //общее количество сообщений
}

func SaveMessage1(db *sql.DB, message *Message) error {
	// Вставляем сообщение в базу данных с использованием NOW()

	log.Printf("Message to be saved: %+v\n", message)

	query := `INSERT INTO msg (text, status_1, created_at_1) VALUES ($1, $2, NOW()) RETURNING id`
	err := db.QueryRow(query, message.Text, message.Status_1).Scan(&message.ID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func SaveMessage2(db *sql.DB, message *Message) error {
	query := `UPDATE msg SET status_2 = $1, created_at_2 = NOW() WHERE id = $2`
	_, err := db.Exec(query, "processed", message.ID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func GetStats(db *sql.DB) (*Stats, error) {
	var stats Stats
	query := `SELECT COUNT(*) FROM msg`
	err := db.QueryRow(query).Scan(&stats.TotalMessages)
	return &stats, err
}
