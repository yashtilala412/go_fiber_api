package config

import (
	"log"

	"github.com/spf13/viper"
)

// LoadConfig initializes Viper
func LoadConfig() {
	viper.SetConfigFile(".env") // ✅ Reads from .env
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// Get returns a config value as a string
func Get(key string) string {
	return viper.GetString(key)
}
