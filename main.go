package main

import (
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/cli"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/config"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/logger"

	"github.com/spf13/viper"
)

func main() {
	cfg, _ := config.LoadConfig()

	logger, err := logger.NewRootLogger(cfg.DEBUG, cfg.IS_DEVELOPMENT)
	if err != nil {
		panic(err)
	}

	mode := viper.GetString("MODE")
	if mode == "api" {
		cli.StartAPI(logger)
	} else {
		cli.StartCLI()
	}
}
