package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return errors.Wrap(err, "failed to set goose dialect")
	}

	// Используем относительный путь к директории с миграциями
	migrationsDir := "internal/db/migrations"

	if err := goose.Up(db, migrationsDir); err != nil {
		return errors.Wrap(err, "failed to run goose migrations")
	}

	return nil
}
