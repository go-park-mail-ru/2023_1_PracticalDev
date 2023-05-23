package config

import (
	"github.com/spf13/viper"
)

func GetGRPCAddr() string {
	return viper.GetString(GrpcConfig.Address) + ":" + viper.GetString(GrpcConfig.Port)
}

func GetConsulAddr() string {
	return viper.GetString(GrpcConfig.ConsulAddr) + ":" + viper.GetString(GrpcConfig.Port)
}

func DefaultRedisConfig() {
	viper.Set(RedisConfig.Host, "redis")
	viper.Set(RedisConfig.Port, "6379")
	viper.Set(RedisConfig.Password, "pickpinpswd")
}

func DefaultPostgresConfig() {
	viper.Set(PostgresConfig.Host, "db")
	viper.Set(PostgresConfig.Port, 5432)
	viper.Set(PostgresConfig.DB, "pickpindb")
	viper.Set(PostgresConfig.User, "pickpin")
	viper.Set(PostgresConfig.Password, "pickpinpswd")
	viper.Set(PostgresConfig.SSLMode, "disable")
}

func DefaultS3BucketConfig() {
	viper.Set(S3BucketConfig.BucketName, "pickpin")
	viper.Set(S3BucketConfig.DefaultRegion, "ru-msk")
	viper.Set(S3BucketConfig.Endpoint, "https://hb.bizmrg.com")
}

func DefaultMongoConfig() {
	viper.Set(MongoConfig.URI, "mongodb://mongo:27017")
	viper.Set(MongoConfig.Username, "pickpin")
	viper.Set(MongoConfig.Password, "pickpinpswd")
}

func DefaultConsulConfig() {
	viper.Set(ConsulConfig.Addr, "consul:8500")
}
