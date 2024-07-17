package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(host, port, user, password, dbname string) (*sql.DB, error) {
	// Формирование строки подключения
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Подключение к базе данных
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Проверка соединения
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
