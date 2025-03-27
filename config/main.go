package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// AllConfig variable of type AppConfig
var AllConfig AppConfig

// AppConfig type AppConfig
type AppConfig struct {
	IsDevelopment  bool   `envconfig:"IS_DEVELOPMENT"`
	Debug          bool   `envconfig:"DEBUG"`
	Host           string `envconfig:"HOST"`
	Port           string `envconfig:"APP_PORT"`
	CSVFilePath    string `envconfig:"CSV_FILE_PATH"`
	ReviewFilePath string `envconfig:"REVIEW_FILE_PATH"`
}

// GetConfig Collects all configs
func GetConfig() AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("warning .env file not found, scanning from OS ENV")
	}

	AllConfig = AppConfig{}

	err = envconfig.Process("APP_PORT", &AllConfig)
	if err != nil {
		log.Fatal(err)
	}

	return AllConfig
}

// GetConfigByName Collects all configs
func GetConfigByName(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv(key)
}

// LoadTestEnv loads environment variables from .env.testing file
func LoadTestEnv() AppConfig {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load(fmt.Sprintf("%s/.env.testing", cwd))
	if err != nil {
		log.Fatal(err)
	}
	return GetConfig()
}
