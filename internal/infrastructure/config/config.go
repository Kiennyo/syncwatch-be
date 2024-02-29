package config

import (
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	HTTP     HTTP
	DB       DB
	Security Security
}

type HTTP struct {
	Port int
}

type DB struct {
	Url string
}

type Security struct {
	JWTSecret string
}

func Load() Config {
	return Config{
		HTTP:     loadHTTPConfig(),
		DB:       loadDBConfig(),
		Security: loadSecurity(),
	}
}

func loadHTTPConfig() HTTP {
	return HTTP{
		Port: getEnvAsInt("PORT", 4000),
	}
}

func loadDBConfig() DB {
	return DB{
		Url: getEnv("DB_URL", ""),
	}
}

func loadSecurity() Security {
	return Security{
		JWTSecret: getEnv("JWT_SECRET", "secret"),
	}
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	slog.Warn("env var not set", "key", key)

	return defaultVal
}
