package db

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	_ "github.com/lib/pq"
)

func New(logger log.Logger) (*sql.DB, error) {
	params := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		"db", "5432", "pickpin", "pickpindb", "pickpinpswd", "disable")

	logger.Info("Connecting to db...")
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
