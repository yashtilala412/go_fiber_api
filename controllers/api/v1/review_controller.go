// controllers/api/v1/review_controller.go

package v1

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"errors"

	"github.com/go-playground/validator/v10"
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
	model := models.NewReviewModel(config)
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
// @Success 201 {object} utils.JSONSuccessResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/v1/reviews [post]

func (rc *ReviewController) AddReview(c *fiber.Ctx) error {
	var review models.Review
	body := c.Body()
	if err := json.Unmarshal(body, &review); err != nil {
		rc.logger.Error("Error parsing review data", zap.Error(err))
		return utils.JSONFail(c, fiber.StatusBadRequest, "Invalid review data")
	}

	// Validate the review struct
	validate := validator.New() // Initialize validator here
	if err := validate.Struct(review); err != nil {
		rc.logger.Error("Validation error", zap.Error(err))
		return utils.JSONFail(c, fiber.StatusBadRequest, utils.ValidatorErrorString(err))
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
// @Success 201 {object} utils.JSONSuccessResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 404 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/v1/review/{name} [delete]
func (rc *ReviewController) DeleteReview(c *fiber.Ctx) error {
	// Get the URL-encoded app name parameter
	encodedAppName := c.Params(constants.ParamAppName)

	// URL decode the app name (important for names with spaces and special characters)
	appName, err := url.QueryUnescape(encodedAppName)
	if err != nil {
		rc.logger.Error(constants.ErrDecodingAppName, zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidAppNameFormat) // Use http.StatusBadRequest
	}

	
	// Call the model's DeleteReview method with the decoded name
	if err := rc.reviewModel.DeleteReview(appName); err != nil {
		rc.logger.Error(constants.ErrDeletingReviews, zap.Error(err))
		if err.Error() == constants.AppNotFoundErrorMessage {
			return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAppNotFound) // Use http.StatusBadRequest
		}
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrDeleteReviews) // Use http.StatusInternalServerError
	}
	return utils.JSONSuccess(c, http.StatusOK, constants.ReviewsDeletedSuccessfully) // Use http.StatusOK
}
