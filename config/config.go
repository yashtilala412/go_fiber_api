package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config holds environment variables
type Config struct {
	Port           string `envconfig:"PORT" validate:"required"`
	CSVFilePath    string `envconfig:"CSV_FILE_PATH" validate:"required"`
	ReviewCSVPath  string `envconfig:"REVIEW_CSV_PATH" validate:"required"`
	IS_DEVELOPMENT bool   `envconfig:"IS_DEVELOPMENT" validate:"required"`
	DEBUG          bool   `envconfig:"DEBUG" validate:"required"`
}

// LoadConfig loads configuration from .env and checks file paths
func LoadConfig() (*Config, error) {
	_ = godotenv.Load() // Load .env file if available

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return nil, err
	}

	// Check if CSVFilePath is set and exists
	if cfg.CSVFilePath == "" {
		log.Fatal("Error: CSV_FILE_PATH is missing!")
	} else if _, err := os.Stat(cfg.CSVFilePath); os.IsNotExist(err) {
		log.Fatalf("Error: CSV_FILE_PATH (%s) does not exist!", cfg.CSVFilePath)
	}

	// Check if ReviewCSVPath is set and exists
	if cfg.ReviewCSVPath == "" {
		log.Fatal("Error: REVIEW_CSV_PATH is missing!")
	} else if _, err := os.Stat(cfg.ReviewCSVPath); os.IsNotExist(err) {
		log.Fatalf("Error: REVIEW_CSV_PATH (%s) does not exist!", cfg.ReviewCSVPath)
	}

	// Log loaded config values
	log.Printf("✅ Config Loaded: %+v\n", cfg)
	return &cfg, nil
}
