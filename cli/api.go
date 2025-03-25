package cli

import (
	"log"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// StartAPI initializes the API
func StartAPI(logger *zap.Logger) {
	app := fiber.New()

	// Register routes
	routes.RegisterRoutes(app, logger)

	port := viper.GetString("PORT")
	log.Println("Starting API on port", port)
	log.Fatal(app.Listen(":" + port))
}
