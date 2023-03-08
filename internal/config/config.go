package config

import "os"

var envVars = map[string]string{
	"PGHOST":            "db",
	"PGPORT":            "5432",
	"POSTGRES_DB":       "pickpindb",
	"POSTGRES_USER":     "pickpin",
	"POSTGRES_PASSWORD": "pickpinpswd",
	"LOG_LEVELS":        "*",
	"POSTGRES_SSL":      "disable",
	"REDIS_PASSWORD":    "pickpinpswd",
	"REDIS_HOST":        "redis",
	"REDIS_PORT":        "6379",
}

func init() {
	for envVar := range envVars {
		if envValue := os.Getenv(envVar); envValue != "" {
			envVars[envVar] = envValue
		}
	}
}

// Get returns configuration parameter named variable.
func Get(variable string) string {
	return envVars[variable]
}
