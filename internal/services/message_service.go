package services

import (
	"database/sql"

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
	err := s.DB.QueryRow(query, message.Text, message.Status).Scan(&message.ID)
	return err
}

func (s *MessageService) UpdateMessageStatus(id int, status string) error {
	query := `UPDATE msg SET status = $1 WHERE id = $2`
	_, err := s.DB.Exec(query, status, id)
	return err
}

func (s *MessageService) GetStats() (*models.Stats, error) {
	var stats models.Stats
	query := `SELECT COUNT(*) FROM msg`
	err := s.DB.QueryRow(query).Scan(&stats.TotalMessages)
	return &stats, err
}
