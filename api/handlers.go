package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
// @Param message body models.Message_Request true "Сообщение"
// @Success 201 {object} models.Message_Request "Сообщение успешно создано"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /message [post]
func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message_Request
	fmt.Println("CreateMessage")

	// Декодирование запроса
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		h.Logger.Errorf("Failed to decode message: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Создание структуры сообщения для сохранения в базу данных
	dbMessage := models.Message{
		Text:        message.Text,
		Status_1:    "pending",         // Установить статус по умолчанию
		CreatedAt_1: time.Now().Unix(), // Установить текущее время
		CreatedAt_2: time.Now().Unix(), // Установить текущее время

	}

	// Сохранение сообщения в базе данных
	if err := models.SaveMessage(h.DB, &dbMessage); err != nil {
		h.Logger.Errorf("Failed to save message: %v", err)
		http.Error(w, fmt.Sprintf("Failed to save message: %v", err), http.StatusInternalServerError)
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
	fmt.Println("GetStats")
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
