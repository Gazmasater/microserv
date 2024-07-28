package services

import (
	"database/sql"
	"fmt"

	"github.com/Gazmasater/internal/models"
)

type MessageService struct {
	DB *sql.DB
}

func NewMessageService(db *sql.DB) *MessageService {
	return &MessageService{DB: db}
}

func (s *MessageService) SaveMessage(message *models.Message) error {
	query := `INSERT INTO msg (text, status) VALUES ($1, $2) RETURNING id`
	err := s.DB.QueryRow(query, message.Text, message.Status_1).Scan(&message.ID)
	return err
}

func (s *MessageService) GetStats() (*models.Stats, error) {
	var stats models.Stats
	query := `SELECT COUNT(*) FROM msg`
	err := s.DB.QueryRow(query).Scan(&stats.TotalMessages)
	return &stats, err
}

// MessageExists проверяет, существует ли сообщение с данным ID в базе данных
func MessageExists(db *sql.DB, id int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM public.msg WHERE id = $1)`
	err := db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке существования сообщения: %v", err)
	}
	return exists, nil
}
