// controllers/api/v1/review_controller.go

package v1

import (
	"net/url"
	"strconv"

	"errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/config"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/constants"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/models"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/utils"
)

type ReviewController struct {
	reviewModel *models.ReviewModel
	logger      *zap.Logger
	config      config.AppConfig
}

// NewReviewController initializes the ReviewController with dependencies.
func NewReviewController(logger *zap.Logger, config config.AppConfig) *ReviewController {
	model := models.NewReviewModel(config) // Corrected line
	return &ReviewController{
		reviewModel: model,
		logger:      logger,
		config:      config,
	}
}

// @Summary List reviews
// @Description Get a list of reviews with filters
// @Tags reviews
// @Accept json
// @Produce json
// @Param appName query string false "App Name" default("DefaultAppName")
// @Param sentiment query string false "Sentiment" default("DefaultSentiment")
// @Param polarityMin query number false "Minimum Polarity" default(-1)
// @Param polarityMax query number false "Maximum Polarity" default(1)
// @Success 200 {array} models.Review
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/v1/reviews [get]

// ... (rest of the ReviewController code remains the same)
// ListReviews handles fetching reviews with filters.
func (rc *ReviewController) ListReviews(c *fiber.Ctx) error {
	// Fetch query parameters with default values from constants
	appName := c.Query(constants.ParamAppName, constants.DefaultAppName)
	sentiment := c.Query(constants.ParamSentiment, constants.DefaultSentiment)

	// Parse polarity with more robust error handling
	polarityMin, err := strconv.ParseFloat(c.Query(constants.ParamPolarityMin, constants.DefaultPolarityMin), 64)
	if err != nil {
		rc.logger.Error("Invalid polarity min",
			zap.String("value", c.Query(constants.ParamPolarityMin)),
			zap.Error(err),
		)
		return utils.JSONFail(c, fiber.StatusBadRequest, constants.ErrorInvalidOffset)

	}

	polarityMax, err := strconv.ParseFloat(c.Query(constants.ParamPolarityMax, constants.DefaultPolarityMax), 64)
	if err != nil {
		rc.logger.Error("Invalid polarity max",
			zap.String("value", c.Query(constants.ParamPolarityMax)),
			zap.Error(err),
		)

		return utils.JSONFail(c, fiber.StatusBadRequest, constants.ErrorInvalidOffset)
	}

	// Log the parsed parameters
	rc.logger.Info("Review Query Parameters",
		zap.String("app_name", appName),
		zap.String("sentiment", sentiment),
		zap.Float64("polarity_min", polarityMin),
		zap.Float64("polarity_max", polarityMax),
	)

	// Fetch reviews from the model, passing the fiber context
	reviews, err := rc.reviewModel.ListReviews(c, appName, sentiment, polarityMin, polarityMax)
	if err != nil {

		return utils.JSONFail(c, fiber.StatusBadRequest, err.Error())
	}
	if len(reviews) == 0 {
		err = errors.New("No apps found with filters")
		return utils.JSONSuccess(c, fiber.StatusOK, map[string]interface{}{
			"reviews": reviews,
			"total":   0,
			"message": "No apps found matching the specified filters",
		})
	}

	return utils.JSONSuccess(c, fiber.StatusOK, reviews)
}

// @Summary Add a new review
// @Description Add a new review to the system
// @Tags reviews
// @Accept json
// @Produce json
// @Param review body models.Review true "Review object to be added"
// @Success 201 {string} string "Review added successfully"
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/v1/reviews [post]
func (rc *ReviewController) AddReview(c *fiber.Ctx) error {
	var review models.Review
	if err := c.BodyParser(&review); err != nil {
		rc.logger.Error("Error parsing review data", zap.Error(err))
		return utils.JSONFail(c, fiber.StatusBadRequest, "Invalid review data")
	}

	if err := rc.reviewModel.AddReview(review); err != nil {
		rc.logger.Error("Error adding review", zap.Error(err))
		return utils.JSONFail(c, fiber.StatusInternalServerError, "Failed to add review")
	}

	return utils.JSONSuccess(c, fiber.StatusCreated, "Review added successfully")
}

// @Summary Delete reviews for an app
// @Description Delete all reviews for a given app name
// @Tags reviews
// @Produce json
// @Param name path string true "App name"
// @Success 200 {string} string "Reviews deleted successfully"
// @Failure 400 {object} utils.JSONResponse
// @Failure 404 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/v1/reviews/:name [delete]
func (rc *ReviewController) DeleteReview(c *fiber.Ctx) error {

	// Get the URL-encoded app name parameter
	encodedAppName := c.Params("name")

	// URL decode the app name (important for names with spaces and special characters)
	appName, err := url.QueryUnescape(encodedAppName)
	if err != nil {
		rc.logger.Error("Error decoding app name", zap.Error(err))
		return utils.JSONFail(c, fiber.StatusBadRequest, "Invalid app name format")
	}

	rc.logger.Info("Deleting reviews for app with name", zap.String("appName", appName))

	// Call the model's DeleteReview method with the decoded name
	if err := rc.reviewModel.DeleteReview(appName); err != nil {
		rc.logger.Error("Error deleting reviews", zap.Error(err))
		if err.Error() == "App not found" {
			return utils.JSONFail(c, fiber.StatusNotFound, "App not found")
		}
		return utils.JSONFail(c, fiber.StatusInternalServerError, "Failed to delete reviews")
	}

	return utils.JSONSuccess(c, fiber.StatusOK, "Reviews deleted successfully")
}
