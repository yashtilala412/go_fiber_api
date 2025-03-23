package main

import (
	"log"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Setup routes
	routes.SetupRoutes(app)

	// Start the server
	port := ":3000"
	log.Println("Server is running on http://localhost" + port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
