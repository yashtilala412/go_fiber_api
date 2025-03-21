package main

import (
	"log"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/middleware"
	// "git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Load middleware (cache)
	middleware.NewCsvCache()

	// // Register Routes
	// routes.SetupRoutes(app)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
