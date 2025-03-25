package handlers

import (
	"log"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/services"
	"github.com/gofiber/fiber/v2"
)

// API Handler to fetch application data
func GetAppsHandler(c *fiber.Ctx) error {
	apps := services.GetApps()
	if len(apps) == 0 {
		log.Println("No data available in cache")
		return c.JSON(fiber.Map{"data": nil})
	}

	log.Println("Fetched Apps:", len(apps)) // ✅ Logging fetched rows
	return c.JSON(fiber.Map{"data": apps})
}
