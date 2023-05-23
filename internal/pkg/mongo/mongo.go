package mongo

import (
	"context"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/config"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

func NewMongoClient(logger *zap.Logger) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().
		ApplyURI(viper.GetString(config.MongoConfig.URI)).
		SetAuth(options.Credential{
			Username: viper.GetString(config.MongoConfig.Username),
			Password: viper.GetString(config.MongoConfig.Password),
		}))
	if err != nil {
		logger.Error("failed to connect to mongo db", zap.Error(err))
		return nil, err
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		logger.Error("failed to ping mongo db", zap.Error(err))
		return nil, err
	}
	logger.Debug("connected to mongo db successfully")
	return client, nil
}
