package config

var GrpcConfig = struct {
	ServiceName string
	ConsulAddr  string
	Address     string
	Port        string
	MessageSize string
}{
	ServiceName: "GRPC_SERVICE_NAME",
	ConsulAddr:  "GRPC_CONSUL_ADDR",
	Address:     "GRPC_ADDRESS",
	Port:        "GRPC_PORT",
	MessageSize: "GRPC_MESSAGE_SIZE",
}

var MetricsConfig = struct {
	Addr string
}{
	Addr: "METRICS_ADDR",
}

var HttpConfig = struct {
	Addr string
}{
	Addr: "HTTP_ADDR",
}

var CSRFConfig = struct {
	Token string
}{
	Token: "CSRF_TOKEN_SECRET",
}

var PostgresConfig = struct {
	Host     string
	Port     string
	DB       string
	User     string
	Password string
	SSLMode  string
}{
	Host:     "PGHOST",
	Port:     "PGPORT",
	DB:       "POSTGRES_DB",
	User:     "POSTGRES_USER",
	Password: "POSTGRES_PASSWORD",
	SSLMode:  "POSTGRES_SSL",
}

var SearchPostgresConfig = struct {
	Host     string
	Port     string
	DB       string
	User     string
	Password string
	SSLMode  string
}{
	Host:     "PGHOST",
	Port:     "PGPORT",
	DB:       "POSTGRES_DB",
	User:     "POSTGRES_USER",
	Password: "POSTGRES_PASSWORD",
	SSLMode:  "POSTGRES_SSL",
}

var RedisConfig = struct {
	Host     string
	Port     string
	Password string
}{
	Host:     "REDIS_HOST",
	Port:     "REDIS_PORT",
	Password: "REDIS_PASSWORD",
}

var S3BucketConfig = struct {
	BucketName    string
	DefaultRegion string
	Endpoint      string
}{
	BucketName:    "S3_BUCKET_NAME",
	DefaultRegion: "S3_DEFAULT_REGION",
	Endpoint:      "S3_ENDPOINT",
}

var MongoConfig = struct {
	URI      string
	Username string
	Password string
}{
	URI:      "MONGO_URI",
	Username: "MONGO_USER",
	Password: "MONGO_PASSWORD",
}

var ConsulConfig = struct {
	Addr string
}{
	Addr: "CONSUL_ADDR",
}
