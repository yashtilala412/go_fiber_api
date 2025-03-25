package main

import (
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/cli"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/config"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	config.LoadConfig()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	mode := viper.GetString("MODE")
	if mode == "api" {
		cli.StartAPI(logger)
	} else {
		cli.StartCLI()
	}
}
