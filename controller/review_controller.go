package controller

import (
	"log"
	"net/http"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/services"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/utils"
	"github.com/gofiber/fiber/v2"
)

// API Handler to fetch application data
func GetReviewController(c *fiber.Ctx) error {
	apps := services.GetReview()
	if len(apps) == 0 {
		log.Println("No data available in cache")
		return utils.JSONSuccess(c, http.StatusOK, nil)
	}

	log.Println("Fetched Apps:", len(apps)) // ✅ Logging fetched rows
	return utils.JSONSuccess(c, http.StatusOK, apps)
}
