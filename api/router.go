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

	// Swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Маршрутизация
	r.Post("/message", handler.CreateMessage)
	r.Get("/stats", handler.GetStats)

	return r
}
