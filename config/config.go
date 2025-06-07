package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

type CassandraConfig struct {
	Hosts       []string
	Keyspace    string
	Timeout     time.Duration
	Retries     int
	PoolSize    int
	Consistency gocql.Consistency
}

func LoadCassandraConfig() CassandraConfig {
	timeout, _ := strconv.Atoi(getEnvOrDefault("CASSANDRA_TIMEOUT_SECONDS", "10"))
	retries, _ := strconv.Atoi(getEnvOrDefault("CASSANDRA_MAX_RETRIES", "3"))
	poolSize, _ := strconv.Atoi(getEnvOrDefault("CASSANDRA_POOL_SIZE", "10"))

	consistencyStr := getEnvOrDefault("CASSANDRA_CONSISTENCY", "QUORUM")
	consistency := gocql.ParseConsistency(consistencyStr)

	return CassandraConfig{
		Hosts:       []string{getEnvOrDefault("CASSANDRA_HOST", "cassandra:9042")},
		Keyspace:    getEnvOrDefault("CASSANDRA_KEYSPACE", "techmarket"),
		Timeout:     time.Duration(timeout) * time.Second,
		Retries:     retries,
		PoolSize:    poolSize,
		Consistency: consistency,
	}
}

type MongoDBConfig struct {
	URI string
}

func LoadMongoDBConfig() MongoDBConfig {
	return MongoDBConfig{
		URI: getEnvOrDefault("MONGODB_URI", "mongodb://root:root_password@localhost:27017/techmarket_db?authSource=admin"),
	}
}

type PostgresConfig struct {
	URI string
}

func LoadPostgresConfig() PostgresConfig {
	return PostgresConfig{
		URI: getEnvOrDefault("POSTGRES_URI", "postgres://techmarket_user:techmarket_password@localhost:5436/techmarket_db"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
