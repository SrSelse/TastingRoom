package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type ServerConfig struct {
	dbHost string
	dbUser string
	dbPass string
	dbName string
	dbPort string

	storageDriver string

	centrifugoApi     string
	centrifugoKey     string
	centrifugoHmacKey string

	httpEndpointPort      string
	httpCorsAllowedOrigin []string
	httpReadTimeout       time.Duration
	httpWriteTimeout      time.Duration

	debug bool

	jwtSecret string
}

const httpDefaultTimeout time.Duration = 20 * time.Second

func parseBoolWithDefault(env string, d bool) bool {
	debug, err := strconv.ParseBool(os.Getenv(env))
	if err != nil {
		return d
	}

	return debug
}

func NewConfigFromEnv() (ServerConfig, error) {
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		return ServerConfig{}, fmt.Errorf("unable to parse environment variable DEBUG: %w", err)
	}

	corsAllowedOriginsString := os.Getenv("HTTP_CORS_ALLOWED_ORIGINS")
	corsAllowedOrigins := strings.Split(corsAllowedOriginsString, ",")

	config := ServerConfig{
		dbUser: os.Getenv("DB_USER"),
		dbPass: os.Getenv("DB_PASS"),
		dbHost: os.Getenv("DB_HOST"),
		dbName: os.Getenv("DB_NAME"),
		dbPort: os.Getenv("DB_PORT"),

		storageDriver: os.Getenv("STORAGE_DRIVER"),

		centrifugoApi:     os.Getenv("CENTRIFUGO_API"),
		centrifugoKey:     os.Getenv("CENTRIFUGO_KEY"),
		centrifugoHmacKey: os.Getenv("CENTRIFUGO_HMAC_KEY"),

		httpEndpointPort:      os.Getenv("HTTP_ENDPOINT_PORT"),
		httpCorsAllowedOrigin: corsAllowedOrigins,
		httpReadTimeout:       httpDefaultTimeout,
		httpWriteTimeout:      httpDefaultTimeout,

		debug: debug,

		jwtSecret: os.Getenv("JWT_SECRET"),
	}

	return config, nil
}
