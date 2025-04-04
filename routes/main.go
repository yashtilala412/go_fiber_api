package routes

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/config"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/constants"
	controller "git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/controllers/api/v1"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/middlewares"
	pMetrics "git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/pkg/prometheus"

	"github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

// Setup initializes routes for the application
func Setup(app *fiber.App, logger *zap.Logger, config config.AppConfig, pMetrics *pMetrics.PrometheusMetrics) error {
	mu.Lock()
	defer mu.Unlock()

	app.Use(middlewares.LogHandler(logger, pMetrics))

	router := app.Group("/api")
	v1 := router.Group("/v1")

	// API Endpoints
	SetupAppRoutes(v1, logger, config)
	SetupReviewRoutes(v1, logger, config)

	return nil
}

// SetupAppRoutes defines the routes for app management
func SetupAppRoutes(v1 fiber.Router, logger *zap.Logger, config config.AppConfig) {
	appController := controller.NewAppController(logger, config)

	appGroup := v1.Group("/apps")
	appGroup.Get("/", appController.ListApps) // Fetch apps with limit, page, and price filter
	appGroup.Post("/", appController.AddApp)  // Add a new app
	appGroup.Delete(fmt.Sprintf("/:%s", constants.ParamAppName), appController.DeleteApp)

}

// SetupreviewRoutes defines the routes for app management
func SetupReviewRoutes(v1 fiber.Router, logger *zap.Logger, config config.AppConfig) {

	reviewController := controller.NewReviewController(logger, config)

	reviewGroup := v1.Group("/review")
	reviewGroup.Get("/", reviewController.ListReviews) // Fetch reviews with filters
	reviewGroup.Post("/", reviewController.AddReview)  //add review with given data
	reviewGroup.Delete(fmt.Sprintf("/:%s", constants.ParamAppName), reviewController.DeleteReview)
}
