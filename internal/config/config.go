package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port string	
	DBURL string
	DBMaxConns int
	DBMinConns int
	JWTSecret string
}

func LoadConfig() (*Config, error) {
	// Load all the env variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using OS env variables")
	}

	// Load configuration from environment variables with defaults
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port: getEnv("PORT", "8080"),
		DBURL: getEnv("DB_URL", "postgres://user:password@localhost:5432/dbname?sslmode=disable"),
		DBMaxConns: getEnvAsInt("DB_MAX_CONNS", 10),
		DBMinConns: getEnvAsInt("DB_MIN_CONNS", 1),
		JWTSecret: getEnv("JWT_SECRET", "your_jwt_secret_key"),
	}, nil
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}