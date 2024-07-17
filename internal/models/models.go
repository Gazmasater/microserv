package models

import (
	"database/sql"
	"errors"
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
	ID        string `json:"id" example:"1"`                                            // Уникальный идентификатор сообщения
	Text      string `json:"text" example:"Hello"`                                      // Содержимое сообщения
	Status    string `json:"status" example:"pending" enums:"processed,pending,failed"` // Статус сообщения
	CreatedAt string `json:"created_at" example:"2024-07-17T08:53:00Z"`                 // Дата создания
}

type Stats struct {
	TotalMessages     int `json:"total_messages"`     // Общее количество сообщений
	ProcessedMessages int `json:"processed_messages"` // Количество обработанных сообщений
	PendingMessages   int `json:"pending_messages"`   // Количество сообщений в ожидании
	FailedMessages    int `json:"failed_messages"`    // Количество сообщений, которые не удалось обработать
}

func SaveMessage(db *sql.DB, message *Message) error {
	query := `INSERT INTO msg (text, status) VALUES ($1, $2) RETURNING id`
	err := db.QueryRow(query, message.Text, message.Status).Scan(&message.ID)
	return err
}

func GetStats(db *sql.DB) (*Stats, error) {
	var stats Stats
	query := `SELECT COUNT(*) FROM msg`
	err := db.QueryRow(query).Scan(&stats.TotalMessages)
	return &stats, err
}

func ValidateMessage(message *Message) error {
	if message.ID == "" {
		return ErrEmptyID
	}
	if message.Text == "" {
		return ErrEmptyText
	}
	if message.Status != StatusProcessed && message.Status != StatusPending && message.Status != StatusFailed {
		return ErrInvalidMessage
	}
	return nil
}
