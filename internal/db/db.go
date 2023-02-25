package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func New(logger *log.Logger) (*sql.DB, error) {
	params := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		"db", "5432", "pickpin", "pickpindb", "pickpinpswd", "disable")

	logger.Println("Connecting to db...")
	db, err := sql.Open("postgres", params)
	if err != nil {
		logger.Printf("Failed to create DB connection, %s", err.Error())
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Printf("Failed to connect to DB, %s", err.Error())
		return nil, err
	}
	logger.Println("DB connection created successfully")

	return db, nil
}
