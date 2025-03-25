package routes

import (
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/handlers"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// RegisterRoutes sets up API routes
func RegisterRoutes(app *fiber.App, logger *zap.Logger) {
	api := app.Group("/api")

	// Middleware to inject logger into the request context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("logger", logger)
		return c.Next()
	})

	// Register the GetApps route
	api.Get("/apps", handlers.GetAppsHandler)
}
