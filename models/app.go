package models

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/config"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/utils"
	"github.com/jszwec/csvutil"
	"go.uber.org/zap"
)

// Global Cache Variables
var (
	appCache []App
	appMutex sync.RWMutex
	appOnce  sync.Once
	validate = validator.New()
)

// App represents the structure of each row in CSV
type App struct {
	Name          string  `csv:"App" validate:"required"`
	Category      string  `csv:"Category" validate:"required"`
	Rating        float64 `csv:"Rating" validate:"gte=0,lte=5"`
	Reviews       int     `csv:"Reviews" validate:"gte=0"`
	Size          string  `csv:"Size" validate:"required"`
	Installs      string  `csv:"Installs" validate:"required"`
	Type          string  `csv:"Type" validate:"required"`
	Price         string  `csv:"Price" validate:"required"`
	ContentRating string  `csv:"Content Rating" validate:"required"`
	Genres        string  `csv:"Genres" validate:"required"`
	LastUpdated   string  `csv:"Last Updated" validate:"required"`
	CurrentVer    string  `csv:"Current Ver" validate:"required"`
	AndroidVer    string  `csv:"Android Ver" validate:"required"`
}

// AppModel contains the logger and config
type AppModel struct {
	logger *zap.Logger
	config config.AppConfig
}

// NewAppModel initializes a new AppModel
func NewAppModel(logger *zap.Logger, config config.AppConfig) *AppModel {
	return &AppModel{
		logger: logger,
		config: config,
	}
}

// loadCache: Loads app data into cache
func (am *AppModel) loadCache() error {
	apps, err := am.ParseApps()
	if err != nil {
		return err
	}

	appCache = apps
	return nil
}

// GetAppsFromCache: Returns data from cache or loads it if expired
func (am *AppModel) GetAppsFromCache() ([]App, error) {
	// First-time cache load
	appOnce.Do(func() {
		_ = am.loadCache()
	})

	appMutex.RLock()
	defer appMutex.RUnlock()
	if len(appCache) == 0 {
		err := am.loadCache()
		if err != nil {
			return nil, err
		}
	}

	return appCache, nil
}

// ParseApps: Reads and parses apps from CSV using csvutil
func (am *AppModel) ParseApps() ([]App, error) {
	if am.config.CSVFilePath == "" {
		return nil, errors.New("CSV file path is not configured")
	}

	var apps []App
	records, err := utils.ReadCSV(am.config.CSVFilePath)
	if err != nil {
		return nil, err
	}
	if err := csvutil.Unmarshal(records, &apps); err != nil {
		return nil, err
	}

	// Post-processing: Remove "$" from Price and format Installs
	for i := range apps {
		cleanedPrice := cleanPriceStr(apps[i].Price)
		priceFloat, err := strconv.ParseFloat(cleanedPrice, 64)
		if err != nil {
			priceFloat = 0 // Default to 0 if conversion fails
		}
		apps[i].Price = strconv.FormatFloat(priceFloat, 'f', 2, 64)
		apps[i].Installs = cleanInstalls(apps[i].Installs)
	}
	return apps, nil
}
func (a *App) ValidateApp() error {

	return validate.Struct(a)
}

// ListAllApps: Returns apps with pagination and filters
func (am *AppModel) ListAllApps(limit int, page int, priceFilter string) ([]string, error) {
	apps, err := am.GetAppsFromCache()
	if err != nil {
		return nil, err
	}

	// Apply filters
	var filteredApps []App
	for _, app := range apps {
		// Apply price filter if provided
		if priceFilter != "" {
			price, err := strconv.ParseFloat(priceFilter, 64)
			if err != nil {
				return nil, fmt.Errorf("Invalid price value")
			}

			// Convert app.Price (string) to float64
			appPrice, err := strconv.ParseFloat(app.Price, 64)
			if err != nil {
				continue // Skip invalid price values
			}

			if appPrice != price {
				continue
			}
		}
		filteredApps = append(filteredApps, app)
	}

	totalApps := len(filteredApps)
	offset := (page - 1) * limit
	if offset >= totalApps {
		return []string{}, nil
	}

	// Apply pagination
	end := min(offset+limit, totalApps)
	var appNames []string
	for _, app := range filteredApps[offset:end] {
		appNames = append(appNames, app.Name)
	}

	return appNames, nil
}

func (am *AppModel) AddAppData(app App) error {
	if err := app.ValidateApp(); err != nil {
		return err
	}
	appMutex.Lock()
	defer appMutex.Unlock()

	file, err := os.OpenFile(am.config.CSVFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		app.Name,
		app.Category,
		strconv.FormatFloat(app.Rating, 'f', 6, 64),
		strconv.Itoa(app.Reviews),
		app.Size,
		app.Installs,
		app.Type,
		app.Price,
		app.ContentRating,
		app.Genres,
		app.LastUpdated,
		app.CurrentVer,
		app.AndroidVer,
	}

	if err := writer.Write(record); err != nil {
		return err
	}

	// Append to the in-memory cache
	appCache = append(appCache, app)

	return nil
}

// Helper function for pagination
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// cleanPrice: Removes "$" and converts to float64
func cleanPriceStr(priceStr string) string {
	priceStr = strings.ReplaceAll(priceStr, "$", "") // Remove dollar sign
	if priceStr == "" {
		return "0" // Default to zero if empty
	}
	return priceStr
}

// cleanInstalls: Removes commas from installs field
func cleanInstalls(installs string) string {
	return strings.ReplaceAll(installs, ",", "")
}
