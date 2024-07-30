package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// NewRouter создает и настраивает новый Chi роутер
func NewRouter(db *sql.DB, logger *zap.SugaredLogger) http.Handler {
	r := chi.NewRouter()

	// Добавление middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Создание обработчика
	handler := NewHandler(db, logger)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // путь к сгенерированной документации
	))

	// Маршрутизация
	r.Post("/message", handler.CreateMessage)
	r.Get("/stats", handler.StatsHandler)

	return r
}
