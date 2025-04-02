package v1

import (
	"encoding/json"
	"errors"
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
func (ac *AppController) AddApp(c *fiber.Ctx) error {
	var app models.App
	body := c.Body()
	if err := json.Unmarshal(body, &app); err != nil {
		ac.logger.Error("Error parsing app data", zap.Error(err))
		return utils.JSONFail(c, fiber.StatusBadRequest, "Invalid app data")
	}
	if err := ac.appModel.AddAppData(app); err != nil {
		ac.logger.Error("Error adding app", zap.Error(err))
		return utils.JSONFail(c, fiber.StatusInternalServerError, "Failed to add app")
	}
	return utils.JSONSuccess(c, fiber.StatusCreated, "App added successfully")
}
