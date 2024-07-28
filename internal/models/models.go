package models

import (
	"database/sql"
	"errors"
	"fmt"
)

const (
	StatusProcessed = "processed" // Сообщение было обработано
	StatusPending   = "pending"   // Сообщение ожидает обработки
	StatusFailed    = "failed"    // Обработка сообщения завершилась неудачей
)

const (
	ErrIDRequired      = "id is required"                                     // Ошибка: отсутствует ID
	ErrTextRequired    = "name is required"                                   // Ошибка: отсутствует имя сообщения
	ErrInvalidStatus   = "status must be 'processed', 'pending', or 'failed'" // Ошибка: неверный статус
	ErrMessageNotFound = "message not found"                                  // Ошибка: сообщение не найдено
	ErrDatabase        = "database operation failed"                          // Ошибка: операция с базой данных не удалась
)

var (
	ErrEmptyID        = errors.New(ErrIDRequired)      // Ошибка: ID пустой
	ErrEmptyText      = errors.New(ErrTextRequired)    // Ошибка: текст пустой
	ErrInvalidMessage = errors.New(ErrInvalidStatus)   // Ошибка: неверное сообщение
	ErrNotFound       = errors.New(ErrMessageNotFound) // Ошибка: сообщение не найдено
	ErrDB             = errors.New(ErrDatabase)        // Ошибка: проблема с базой данных
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
	Text string `json:"text"`
}

type Stats struct {
	TotalMessages     int `json:"total_messages"`     // Общее количество сообщений
	ProcessedMessages int `json:"processed_messages"` // Количество обработанных сообщений
	PendingMessages   int `json:"pending_messages"`   // Количество сообщений в ожидании
	FailedMessages    int `json:"failed_messages"`    // Количество сообщений, которые не удалось обработать
}

func SaveMessage(db *sql.DB, message *Message) error {
	// Вставляем сообщение в базу данных с использованием NOW()
	query := `INSERT INTO msg (text, status_1) VALUES ($1, $2) RETURNING id`
	err := db.QueryRow(query, message.Text, message.Status_1).Scan(&message.ID)
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
func ValidateMessage(message *Message) error {
	if message.ID == 0 {
		return ErrEmptyID
	}
	if message.Text == "" {
		return ErrEmptyText
	}

	return nil
}
