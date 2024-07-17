package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB) error {

	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		return errors.Wrap(err, "failed to set goose dialect")
	}

	if err := goose.Up(db, "internal/db/migrations"); err != nil {
		return errors.Wrap(err, "failed to run goose migrations")
	}

	return nil
}
