package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/config"
)

func New(logger *zap.Logger) (*sql.DB, error) {
	logger.Info("Connecting to db...",
		zap.String("host", config.Get("PGHOST")),
		zap.String("port", config.Get("PGPORT")),
		zap.String("dbname", config.Get("POSTGRES_DB")),
		zap.String("sslmode", config.Get("POSTGRES_SSL")))

	params := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		config.Get("PGHOST"), config.Get("PGPORT"), config.Get("POSTGRES_USER"),
		config.Get("POSTGRES_DB"), config.Get("POSTGRES_PASSWORD"), config.Get("POSTGRES_SSL"))
	db, err := sql.Open("postgres", params)
	if err != nil {
		logger.Error("Failed to create DB connection", zap.Error(err))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Error("Failed to connect to DB", zap.Error(err))
		return nil, err
	}
	logger.Info("DB connection created successfully")

	return db, nil
}
