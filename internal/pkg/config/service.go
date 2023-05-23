package config

import "github.com/spf13/viper"

func setGRPCServiceConfig(ServiceName string, Address string, Port int, MessageSize int, ConsulAddr string) {
	viper.Set(GrpcConfig.ServiceName, ServiceName)
	viper.Set(GrpcConfig.Address, Address)
	viper.Set(GrpcConfig.Port, Port)
	viper.Set(GrpcConfig.MessageSize, MessageSize)
	viper.Set(GrpcConfig.ConsulAddr, ConsulAddr)
}

func setupMetricsConfig(Addr string) {
	viper.Set(MetricsConfig.Addr, "0.0.0.0:9003")
}

func setupHTTPConfig(Addr string) {
	viper.Set(HttpConfig.Addr, Addr)
}

func setupCSRFSecretToken(token string) {
	viper.Set(CSRFConfig.Token, token)
}

func DefaultGRPCAuthConfig() {
	setGRPCServiceConfig("auth", "0.0.0.0", 8087, 10, "auth")
	setupMetricsConfig("0.0.0.0:9003")

	DefaultPostgresConfig()
	DefaultRedisConfig()
	DefaultConsulConfig()
}

func DefaultGRPCImageConfig() {
	setGRPCServiceConfig("image", "0.0.0.0", 8088, 10, "image")
	setupMetricsConfig("0.0.0.0:9002")

	DefaultS3BucketConfig()
	DefaultConsulConfig()
}

func DefaultGRPCSearchConfig() {
	setGRPCServiceConfig("search", "0.0.0.0", 8089, 10, "search")
	setupMetricsConfig("0.0.0.0:9004")

	DefaultPostgresConfig()
	DefaultConsulConfig()
}

func DefaultGRPCShortenerConfig() {
	setGRPCServiceConfig("shortener", "0.0.0.0", 8090, 10, "shortener")
	setupMetricsConfig("0.0.0.0:9005")
	setupHTTPConfig("0.0.0.0:8091")

	DefaultMongoConfig()
	DefaultConsulConfig()
}

func DefaultPickPinConfig() {
	setupMetricsConfig("0.0.0.0:9001")
	setupHTTPConfig("0.0.0.0:8080")
	setupCSRFSecretToken("pickpinsecret")

	DefaultPostgresConfig()
	DefaultConsulConfig()
}
