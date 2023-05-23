package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/config"
)

func New(logger *zap.Logger) (*sql.DB, error) {
	logger.Info("Connecting to db...",
		zap.String("host", viper.GetString(config.PostgresConfig.Host)),
		zap.Int("port", viper.GetInt(config.PostgresConfig.Port)),
		zap.String("dbname", viper.GetString(config.PostgresConfig.DB)),
		zap.String("sslmode", viper.GetString(config.PostgresConfig.SSLMode)))

	params := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		viper.GetString(config.PostgresConfig.Host), viper.GetInt(config.PostgresConfig.Port), viper.GetString(config.PostgresConfig.User),
		viper.GetString(config.PostgresConfig.DB), viper.GetString(config.PostgresConfig.Password), viper.GetString(config.PostgresConfig.SSLMode))
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
