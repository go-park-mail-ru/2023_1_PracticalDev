package redis

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/log"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(logger log.Logger, ctx context.Context) (*redis.Client, error) {
	logger.Info("Connecting to redis...")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "pickpinpswd",
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to create Redis connection, ", err.Error())
	}

	logger.Info("Redis connection created successfully")
	return rdb, nil
}
