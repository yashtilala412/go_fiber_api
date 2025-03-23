package handlers

import (
	"log"
	"strconv"

	service "git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/services"
	"github.com/gofiber/fiber/v2"
)

// GetAllApps fetches all apps with pagination and filtering
func GetAllApps(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "10")) // Default limit is 10
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit parameter"})
	}

	page, err := strconv.Atoi(c.Query("page", "1")) // Default page is 1
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page parameter"})
	}

	price := c.Query("price", "")

	apps, err := service.FetchAllApps(limit, page, price)
	if err != nil {
		log.Println("Error fetching apps:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch apps"})
	}

	return c.JSON(fiber.Map{"apps": apps})
}
