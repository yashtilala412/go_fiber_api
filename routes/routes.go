package routes

import (
	handlers "git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/handler"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the API endpoints
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Route to get all apps
	v1.Get("/apps", handlers.GetAllApps)
}
