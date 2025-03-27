package v1

import (
	"strconv"

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
	model := models.NewReviewModel(logger, config)
	return &ReviewController{
		reviewModel: model,
		logger:      logger,
		config:      config,
	}
}

// GetReviews handles fetching reviews with filters.
func (rc *ReviewController) GetReviews(c *fiber.Ctx) error {
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
		polarityMin = 0.0 // Set a safe default
	}

	polarityMax, err := strconv.ParseFloat(c.Query(constants.ParamPolarityMax, constants.DefaultPolarityMax), 64)
	if err != nil {
		rc.logger.Error("Invalid polarity max",
			zap.String("value", c.Query(constants.ParamPolarityMax)),
			zap.Error(err),
		)
		polarityMax = 1.0 // Set a safe default
	}

	// Log the parsed parameters
	rc.logger.Info("Review Query Parameters",
		zap.String("app_name", appName),
		zap.String("sentiment", sentiment),
		zap.Float64("polarity_min", polarityMin),
		zap.Float64("polarity_max", polarityMax),
	)

	// Fetch reviews from the model
	reviews, err := rc.reviewModel.GetReviews(appName, sentiment, polarityMin, polarityMax)
	if err != nil {
		rc.logger.Error("Error fetching reviews", zap.Error(err))
		return utils.JSONFail(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.JSONSuccess(c, fiber.StatusOK, reviews)
}
