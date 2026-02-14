package config

import (
	"os"
	"strconv"
)

type Config struct {
	PostgresDSN string
	RedisAddr   string
	HTTPPort    int

	LoggingLevel string

	KafkaBrokers            string
	KafkaGroupID            string
	KafkaTimeoutMS          int
	KafkaTopic              string
	KafkaReadMessageTimeout int
	KafkaAutoOffsetStore    bool
	KafkaAutoCommit         bool
	KafkaAutoCommitInterval int
}

func Load() *Config {
	return &Config{
		PostgresDSN: getEnv("POSTGRES_DSN", "postgres://user:pass@localhost:5432/db?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		HTTPPort:    getEnvAsInt("HTTP_PORT", 8080),

		KafkaBrokers:            getEnv("KAFKA_BROKERS", ""),
		KafkaGroupID:            getEnv("KAFKA_GROUP_ID", ""),
		KafkaTimeoutMS:          getEnvAsInt("KAFKA_TIMEOUT_MS", -1),
		KafkaTopic:              getEnv("KAFKA_TOPIC", ""),
		KafkaReadMessageTimeout: getEnvAsInt("KAFKA_READ_MESSAGE_TIMEOUT", -1),
		KafkaAutoOffsetStore:    getEnvAsBool("KAFKA_AUTO_OFFSET_STORE", false),
		KafkaAutoCommit:         getEnvAsBool("KAFKA_AUTO_COMMIT", true),
		KafkaAutoCommitInterval: getEnvAsInt("KAFKA_AUTO_COMMIT_INTERVAL", -1),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	switch valueStr {
	case "true":
		return true
	case "false":
		return false
	default:
		return defaultValue
	}
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultValue
}
