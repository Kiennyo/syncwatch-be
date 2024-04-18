package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTP     HTTP
	DB       DB
	Security Security
	SMTP     SMPT
}

type HTTP struct {
	Port int
}

type DB struct {
	URL         string
	MaxOpenConn int
	MaxIdleConn int
	MaxIdleTime string
}

type Security struct {
	JWTSecret string
	Iss       string
	Aud       string
}

type SMPT struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file, defaulting to environment variables")
	}

	config := Config{
		HTTP:     loadHTTPConfig(),
		DB:       loadDBConfig(),
		Security: loadSecurity(),
		SMTP:     loadSMTPConfig(),
	}

	flag.Parse()

	return config
}

func loadHTTPConfig() HTTP {
	http := HTTP{}
	setEnvInt(&http.Port, "PORT", "API server port")

	return http
}

func loadDBConfig() DB {
	db := DB{}
	setEnv(&db.URL, "DB_URL", "PostgreSQL DSN")
	setEnvInt(&db.MaxOpenConn, "DB_MAX_OPEN_CONN", "PostgreSQL max open connections")
	setEnvInt(&db.MaxIdleConn, "DB_MAX_OPEN_IDLE", "PostgreSQL max idle connections")
	setEnv(&db.MaxIdleTime, "DB_MAX_IDLE_TIME", "PostgreSQL max connection idle time (mins)")

	return db
}

func loadSecurity() Security {
	security := Security{}
	setEnv(&security.JWTSecret, "JWT_SECRET", "Secret key to create and verify JWT")
	setEnv(&security.Iss, "JWT_ISS", "JWT issuer")
	setEnv(&security.Aud, "JWT_AUD", "JWT audience")

	return security
}

func loadSMTPConfig() SMPT {
	smpt := SMPT{}
	setEnv(&smpt.Host, "SMPT_HOST", "SMTP host")
	setEnvInt(&smpt.Port, "SMPT_PORT", "SMTP port")
	setEnv(&smpt.Username, "SMPT_USERNAME", "SMTP username")
	setEnv(&smpt.Password, "SMPT_PASSWORD", "SMTP password")
	setEnv(&smpt.Sender, "SMPT_SENDER", "SMTP sender")

	return smpt
}

func setEnvInt(configValue *int, key string, usage string) {
	if envValue, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(envValue); err == nil {
			flag.IntVar(configValue, key, value, usage)
			return
		}
	}

	panic(fmt.Sprintf("env var: %s, can't set or cannot be converted to number", key))
}

func setEnv(configValue *string, key string, usage string) {
	if envValue, exists := os.LookupEnv(key); exists {
		flag.StringVar(configValue, key, envValue, usage)
		return
	}

	panic(fmt.Sprintf("env var: %s, not set", key))
}
