package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	ENS    ENSConfig
	Cache  CacheConfig
	Redis  RedisConfig
	API    APIConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type ENSConfig struct {
	ProviderURL    string
	TimeoutSeconds int
	RetryAttempts  int
}

type CacheConfig struct {
	Type string
	TTL  time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type APIConfig struct {
	EnableRateLimit bool
	RateLimit       struct {
		Requests int
		Duration time.Duration
	}
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{}

	// Server Config
	cfg.Server.Host = getEnvString("SERVER_HOST", "localhost")
	port, err := getEnvInt("SERVER_PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}
	cfg.Server.Port = port

	// ENS Config
	cfg.ENS.ProviderURL = getEnvString("ENS_PROVIDER_URL", "")
	if cfg.ENS.ProviderURL == "" {
		return nil, fmt.Errorf("ENS_PROVIDER_URL is required")
	}
	timeoutSecs, err := getEnvInt("ENS_TIMEOUT_SECONDS", 10)
	if err != nil {
		return nil, fmt.Errorf("invalid ENS_TIMEOUT_SECONDS: %w", err)
	}
	cfg.ENS.TimeoutSeconds = timeoutSecs

	retryAttempts, err := getEnvInt("ENS_RETRY_ATTEMPTS", 3)
	if err != nil {
		return nil, fmt.Errorf("invalid ENS_RETRY_ATTEMPTS: %w", err)
	}
	cfg.ENS.RetryAttempts = retryAttempts

	// Cache Config
	cfg.Cache.Type = getEnvString("CACHE_TYPE", "memory")
	ttlMinutes, err := getEnvInt("CACHE_TTL_MINUTES", 60)
	if err != nil {
		return nil, fmt.Errorf("invalid CACHE_TTL_MINUTES: %w", err)
	}
	cfg.Cache.TTL = time.Duration(ttlMinutes) * time.Minute

	// Redis Config
	cfg.Redis.Host = getEnvString("REDIS_HOST", "localhost")
	redisPort, err := getEnvInt("REDIS_PORT", 6379)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_PORT: %w", err)
	}
	cfg.Redis.Port = redisPort
	cfg.Redis.Password = getEnvString("REDIS_PASSWORD", "")
	redisDB, err := getEnvInt("REDIS_DB", 0)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
	}
	cfg.Redis.DB = redisDB

	// API Config
	cfg.API.EnableRateLimit = getEnvBool("ENABLE_RATE_LIMIT", true)
	requests, err := getEnvInt("RATE_LIMIT_REQUESTS", 100)
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_REQUESTS: %w", err)
	}
	cfg.API.RateLimit.Requests = requests

	durationSeconds, err := getEnvInt("RATE_LIMIT_DURATION_SECONDS", 60)
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_DURATION_SECONDS: %w", err)
	}
	cfg.API.RateLimit.Duration = time.Duration(durationSeconds) * time.Second

	return cfg, nil
}

func getEnvString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue, nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue, fmt.Errorf("invalid value for %s: %w", key, err)
	}
	return intValue, nil
}

func getEnvBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}
