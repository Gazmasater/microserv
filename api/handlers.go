package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Gazmasater/kafka"

	"github.com/Gazmasater/internal/models"
	"go.uber.org/zap"
)

// Handler структура для обработки запросов
type Handler struct {
	DB     *sql.DB
	Logger *zap.SugaredLogger
}

// NewHandler создает новый обработчик
func NewHandler(db *sql.DB, logger *zap.SugaredLogger) *Handler {
	return &Handler{DB: db, Logger: logger}
}

// CreateMessage создает новое сообщение
// @Summary Создать сообщение
// @Description Создает новое сообщение и сохраняет его в базе данных
// @Tags messages
// @Accept json
// @Produce json
// @Param message body models.Message true "Сообщение"
// @Success 201 {object} models.Message "Сообщение успешно создано"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /message [post]
func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	fmt.Println("CreateMessage")
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		h.Logger.Errorf("Failed to decode message: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Println("CreateMessage message", message)

	// Валидация сообщения
	if err := models.ValidateMessage(&message); err != nil {
		h.Logger.Errorf("Validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("CreateMessage перед сохр в базу")

	// Сохранение сообщения в базе данных
	if err := models.SaveMessage(h.DB, &message); err != nil {
		h.Logger.Errorf("Failed to save message: %v", err)
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	// Отправка сообщения в Kafka
	if err := kafka.ProduceMessage(message); err != nil {
		h.Logger.Errorf("Failed to produce message to Kafka: %v", err)
		http.Error(w, "Failed to produce message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetStats обрабатывает GET-запрос для получения статистики
// @Summary Get statistics
// @Description Get statistics from the database
// @Tags stats
// @Produce json
// @Param limit query int false "Limit the number of results"
// @Success 200 {object} models.Stats "Statistics retrieved successfully"
// @Failure 500 {string} string "Failed to get or encode stats"
// @Router /stats [get]
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := models.GetStats(h.DB)
	if err != nil {
		h.Logger.Errorf("Failed to get stats: %v", err)
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(stats); err != nil {
		h.Logger.Errorf("Failed to encode stats: %v", err)
		http.Error(w, "Failed to encode stats", http.StatusInternalServerError)
		return
	}
}

// SendMessage отправляет сообщение в Kafka
// @Summary Отправить сообщение в Kafka
// @Description Отправляет сообщение, полученное в теле запроса, в Kafka.
// @Tags messages
// @Accept json
// @Produce json
// @Param message body models.Message true "Сообщение для отправки"
// @Success 200 {string} string "Message sent"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Failed to produce message"
// @Router /send [post]
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		h.Logger.Errorf("Error decoding message: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := kafka.ProduceMessage(message); err != nil {
		h.Logger.Errorf("Failed to produce message: %v", err)
		http.Error(w, "Failed to produce message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message sent"))
}
