// controllers/api/v1/review_controller.go

package v1

import (
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
