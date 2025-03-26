package routes

import (
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/controller"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// RegisterRoutes sets up API routes
func RegisterRoutes(app *fiber.App, logger *zap.Logger) {
	api := app.Group("/api")

	// Register the GetApps route
	api.Get("/apps", controller.GetAppsController)
	api.Get("/review", controller.GetReviewController)
}
