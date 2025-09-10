package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	App      AppConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
	Issuer     string
}

type AppConfig struct {
	ServingsPerCredit   int
	TokenMultiplier     int
	ClaimCooloffSeconds int
	RateLimitRPS        int
	RateLimitBurst      int
	IdempotencyTTL      time.Duration
	MonthlyCreditsBase  int
	MonthlyCreditsK     float64
	RemoteProxyBonus    float64
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", "30s"),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", "30s"),
			IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", "120s"),
		},
		Database: DatabaseConfig{
			URL:             os.Getenv("DB_URL"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", "5m"),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key"),
			Expiration: getDurationEnv("JWT_EXPIRATION", "24h"),
			Issuer:     getEnv("JWT_ISSUER", "foodflow"),
		},
		App: AppConfig{
			ServingsPerCredit:   getIntEnv("SERVINGS_PER_CREDIT", 10),
			TokenMultiplier:     getIntEnv("TOKEN_MULTIPLIER", 5),
			ClaimCooloffSeconds: getIntEnv("CLAIM_COOLOFF_SECONDS", 600),
			RateLimitRPS:        getIntEnv("RATE_LIMIT_RPS", 10),
			RateLimitBurst:      getIntEnv("RATE_LIMIT_BURST", 20),
			IdempotencyTTL:      getDurationEnv("IDEMPOTENCY_TTL", "24h"),
			MonthlyCreditsBase:  getIntEnv("MONTHLY_CREDITS_BASE", 50),
			MonthlyCreditsK:     getFloatEnv("MONTHLY_CREDITS_K", 0.2),
			RemoteProxyBonus:    getFloatEnv("REMOTE_PROXY_BONUS", 0.1),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getFloatEnv(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getDurationEnv(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}
