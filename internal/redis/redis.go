package redis

import (
	"context"

	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/config"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedisClient(logger *zap.Logger, ctx context.Context) (*redis.Client, error) {
	logger.Info("Connecting to redis...")

	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString(config.RedisConfig.Host) + ":" + viper.GetString(config.RedisConfig.Port),
		Password: viper.GetString(config.RedisConfig.Password),
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to create Redis connection, ", zap.Error(err))
	}

	logger.Info("Redis connection created successfully")
	return rdb, nil
}
