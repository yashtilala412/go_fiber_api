package cli

import (
	"os"
	"os/signal"
	"syscall"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/config"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/routes"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// StartAPI initializes the API
func StartAPI(logger *zap.Logger) {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	app := fiber.New()

	// Load CSV Data
	_, err = utils.LoadCSV(logger, cfg.CSVFilePath) // Use the correct field for CSV path
	if err != nil {
		logger.Fatal("Failed to load CSV data", zap.Error(err))
	}

	// Register routes
	routes.RegisterRoutes(app, logger)

	// Handle graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			logger.Panic("Error starting server", zap.Error(err))
		}
	}()

	<-interrupt
	logger.Info("Gracefully shutting down...")

	if err := app.Shutdown(); err != nil {
		logger.Panic("Error while shutting down server", zap.Error(err))
	}

	logger.Info("Server stopped accepting new requests.")
}
