package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	Redis         RedisConfig
	OpenTelemetry OpenTelemetryConfig
	JWT           JWTConfig
	Elasticsearch ElasticsearchConfig
	MinIO         MinIOConfig
	RabbitMQ      RabbitMQConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         int
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// OpenTelemetryConfig holds OpenTelemetry configuration
type OpenTelemetryConfig struct {
	Endpoint       string
	ServiceName    string
	ServiceVersion string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret            string
	Expiration        time.Duration
	RefreshExpiration time.Duration
}

// ElasticsearchConfig holds Elasticsearch configuration
type ElasticsearchConfig struct {
	URL      string
	Username string
	Password string
}

// MinIOConfig holds MinIO configuration
type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

// RabbitMQConfig holds RabbitMQ configuration
type RabbitMQConfig struct {
	URL      string
	Username string
	Password string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (optional, error is ignored)
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DATABASE_HOST", "localhost"),
			Port:            getEnvAsInt("DATABASE_PORT", 5432),
			User:            getEnv("DATABASE_USER", "postgres"),
			Password:        getEnv("DATABASE_PASSWORD", "password"),
			DBName:          getEnv("DATABASE_NAME", "fiscaflow"),
			SSLMode:         getEnv("DATABASE_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DATABASE_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DATABASE_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DATABASE_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		OpenTelemetry: OpenTelemetryConfig{
			Endpoint:       getEnv("OTEL_ENDPOINT", "jaeger:4317"),
			ServiceName:    getEnv("OTEL_SERVICE_NAME", "fiscaflow"),
			ServiceVersion: getEnv("OTEL_SERVICE_VERSION", "1.0.0"),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "your-secret-key"),
			Expiration:        getEnvAsDuration("JWT_EXPIRATION", 15*time.Minute),
			RefreshExpiration: getEnvAsDuration("JWT_REFRESH_EXPIRATION", 7*24*time.Hour),
		},
		Elasticsearch: ElasticsearchConfig{
			URL:      getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
			Username: getEnv("ELASTICSEARCH_USERNAME", ""),
			Password: getEnv("ELASTICSEARCH_PASSWORD", ""),
		},
		MinIO: MinIOConfig{
			Endpoint:        getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKeyID:     getEnv("MINIO_ACCESS_KEY_ID", "minioadmin"),
			SecretAccessKey: getEnv("MINIO_SECRET_ACCESS_KEY", "minioadmin"),
			UseSSL:          getEnvAsBool("MINIO_USE_SSL", false),
		},
		RabbitMQ: RabbitMQConfig{
			URL:      getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			Username: getEnv("RABBITMQ_USERNAME", "guest"),
			Password: getEnv("RABBITMQ_PASSWORD", "guest"),
		},
	}

	return config, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// GetRedisAddr returns the Redis address
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
