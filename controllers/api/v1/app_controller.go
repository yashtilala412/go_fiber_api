package v1

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/config"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/constants"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/models"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/utils"
)

type AppController struct {
	appModel *models.AppModel
	logger   *zap.Logger
	config   config.AppConfig
}

// NewAppController initializes the AppController with dependencies.
func NewAppController(logger *zap.Logger, config config.AppConfig) *AppController {
	model := models.NewAppModel(logger, config)
	return &AppController{
		appModel: model,
		logger:   logger,
		config:   config,
	}
}

// @Summary List apps
// @Description Get a list of apps with pagination and filters
// @Tags apps
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param priceFilter query string false "Price Filter"
// @Success 200 {array} models.App
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/v1/apps [get]

// ListApps handles the request for fetching all apps with pagination and filters.
func (ac *AppController) ListApps(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query(constants.Limit, constants.DefaultLimit))
	if err != nil || limit <= 0 {
		ac.logger.Error(constants.ErrorInvalidLimit+err.Error(), zap.Error(err))
		return utils.JSONError(c, fiber.StatusBadRequest, constants.ErrorInvalidLimit)
	}

	offset, err := strconv.Atoi(c.Query(constants.Offset, constants.DefaultOffset))
	if err != nil || offset < 0 {
		ac.logger.Error(constants.ErrorInvalidOffset+err.Error(), zap.Error(err))
		return utils.JSONFail(c, fiber.StatusBadRequest, constants.ErrorInvalidOffset)
	}

	// Extract filters from query parameters
	priceFilter := c.Query(constants.ParamFilterPrice, "")

	apps, err := ac.appModel.ListAllApps(limit, offset, priceFilter)
	if err != nil {
		return utils.JSONFail(c, fiber.StatusBadRequest, err.Error())
	}

	if len(apps) == 0 {
		err = errors.New("no apps found with filters")
		ac.logger.Error(err.Error(), zap.Error(err))
		return utils.JSONSuccess(c, fiber.StatusOK, map[string]interface{}{
			"apps":    apps,
			"total":   0,
			"message": "No apps found matching the specified filters",
		})
	}

	return utils.JSONSuccess(c, fiber.StatusOK, apps)
}

// @Summary Add a new app
// @Description Add a new app to the system
// @Tags apps
// @Accept json
// @Produce json
// @Param app body models.App true "App object to be added"
// @Success 201 {string} string "App added successfully"
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/v1/apps [post]
func (ac *AppController) AddApp(c *fiber.Ctx) error {

	var app models.App
	if err := c.BodyParser(&app); err != nil {
		ac.logger.Error("Error parsing app data", zap.Error(err))
		return utils.JSONFail(c, fiber.StatusBadRequest, "Invalid app data")
	}

	if err := ac.appModel.AddAppData(app); err != nil {
		ac.logger.Error("Error adding app", zap.Error(err))
		return utils.JSONFail(c, fiber.StatusInternalServerError, "Failed to add app")
	}

	return utils.JSONSuccess(c, fiber.StatusCreated, "App added successfully")
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
// @Router /api/v1/apps/:name [delete]
func (ac *AppController) DeleteApp(c *fiber.Ctx) error {

	// Get the URL-encoded app name parameter
	encodedAppName := c.Params("name")

	// URL decode the app name (important for names with spaces and special characters)
	appName, err := url.QueryUnescape(encodedAppName)
	if err != nil {
		ac.logger.Error("Error decoding app name", zap.Error(err))
		return utils.JSONFail(c, fiber.StatusBadRequest, "Invalid app name format")
	}

	ac.logger.Info("Deleting app with name", zap.String("appName", appName))

	// Call the model's DeleteApp method with the decoded name
	if err := ac.appModel.DeleteApp(appName); err != nil {
		ac.logger.Error("Error deleting app", zap.Error(err))
		if err.Error() == "App not found" {
			return utils.JSONFail(c, fiber.StatusNotFound, "App not found")
		}
		return utils.JSONFail(c, fiber.StatusInternalServerError, "Failed to delete app")
	}

	return utils.JSONSuccess(c, fiber.StatusOK, "App deleted successfully")
}
