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

	// Проверка на пустую строку
	if message.Text == "" {
		h.Logger.Errorf("Message text is empty")
		http.Error(w, "Message text cannot be empty", http.StatusBadRequest)
		return
	}

	// Преобразование строки в указатель на строку
	text := message.Text

	// Создание структуры сообщения для сохранения в базу данных
	dbMessage := models.Message{
		ID:          0,
		Text:        text,
		Status_1:    "pending",         // Установить статус по умолчанию
		CreatedAt_1: time.Now().Unix(), // Установить текущее время
		CreatedAt_2: time.Now().Unix(), // Установить текущее время
	}

	// Сохранение сообщения в базе данных
	if err := models.SaveMessage1(h.DB, &dbMessage); err != nil {
		h.Logger.Errorf("Failed to save message: %v", err)
		http.Error(w, fmt.Sprintf("Failed to save message: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// StatsHandler обрабатывает GET-запрос для получения статистики
// @Summary Get statistics
// @Description Get statistics from the database
// @Tags stats
// @Produce json
// @Param limit query int false "Limit the number of results"
// @Success 200 {object} models.Stats "Statistics retrieved successfully"
// @Failure 500 {string} string "Failed to get or encode stats"
// @Router /stats [get]
func (h *Handler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	// Параметры подключения к базе данных
	connStr := "user=postgres password=qwert dbname=microserv host=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	limit := r.URL.Query().Get("limit")
	var stats *models.Stats // Измените на указатель
	if limit != "" {
		stats, err = models.GetStatsWithLimit(db, limit)
	} else {
		stats, err = models.GetStats(db) // здесь также необходимо будет изменить
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetStats(db *sql.DB, limit int) (*models.Stats, error) {
	var stats models.Stats
	query := `
		SELECT
			COUNT(*) FILTER (WHERE status_1 = 'pending') AS pending_messages,
			COUNT(*) FILTER (WHERE status_2 = 'processed') AS processed_messages,
			COUNT(*) AS total_messages
		FROM (SELECT * FROM msg LIMIT $1) AS limited_msg;
	`
	err := db.QueryRow(query, limit).Scan(&stats.PendingMessages, &stats.ProcessedMessages, &stats.TotalMessages)
	return &stats, err
}
