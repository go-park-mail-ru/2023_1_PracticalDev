package db

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/config"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	_ "github.com/lib/pq"
)

func New(logger log.Logger) (*sql.DB, error) {
	params := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		config.Get("PGHOST"), config.Get("PGPORT"), config.Get("POSTGRES_USER"), config.Get("POSTGRES_DB"), config.Get("POSTGRES_PASSWORD"), config.Get("POSTGRES_SSL"))

	logger.Info("Connecting to db...", params)
	db, err := sql.Open("postgres", params)
	if err != nil {
		logger.Error("Failed to create DB connection, ", err.Error())
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Error("Failed to connect to DB, ", err.Error())
		return nil, err
	}
	logger.Info("DB connection created successfully")

	return db, nil
}
