package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort int
}

// LoadConfig reads configuration from environment variables or file.
func LoadConfig() (config Config, err error) {
	// Load .env file
	err = godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// Set defaults
	config.DBHost = "localhost"
	config.DBPort = 5432
	config.DBUser = "postgres"
	config.DBPassword = ""
	config.DBName = "go_todo_db"
	config.ServerPort = 8080

	// Override from environment variables if present
	if os.Getenv("DB_HOST") != "" {
		config.DBHost = os.Getenv("DB_HOST")
	}

	if os.Getenv("DB_PORT") != "" {
		config.DBPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	}

	if os.Getenv("DB_USER") != "" {
		config.DBUser = os.Getenv("DB_USER")
	}

	if os.Getenv("DB_PASSWORD") != "" {
		config.DBPassword = os.Getenv("DB_PASSWORD")
	}

	if os.Getenv("DB_NAME") != "" {
		config.DBName = os.Getenv("DB_NAME")
	}

	if os.Getenv("SERVER_PORT") != "" {
		config.ServerPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
	}

	return config, nil
}

// GetDBConnString returns the database connection string
func (c Config) GetDBConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}
