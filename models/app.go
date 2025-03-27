package models

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"strings"
	"sync"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/config"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jszwec/csvutil"
	"go.uber.org/zap"
)

// App represents the structure of each row in CSV
type App struct {
	Name          string  `csv:"App"`
	Category      string  `csv:"Category"`
	Rating        float64 `csv:"Rating"`
	Reviews       int     `csv:"Reviews"`
	Size          string  `csv:"Size"`
	Installs      string  `csv:"Installs"`
	Type          string  `csv:"Type"`
	Price         string  `csv:"Price"`
	ContentRating string  `csv:"Content Rating"`
	Genres        string  `csv:"Genres"`
	LastUpdated   string  `csv:"Last Updated"`
	CurrentVer    string  `csv:"Current Ver"`
	AndroidVer    string  `csv:"Android Ver"`
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

// Global Cache Variables
var (
	appCache []App
	appMutex sync.RWMutex
	appOnce  sync.Once
)

// loadCache: Loads app data into cache
func (am *AppModel) loadCache() error {
	appMutex.Lock()
	defer appMutex.Unlock()

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
	if len(appCache) == 0 {
		appMutex.RUnlock()
		err := am.loadCache()
		if err != nil {
			return nil, err
		}
		appMutex.RLock()
	}
	defer appMutex.RUnlock()

	return appCache, nil
}

// ParseApps: Reads and parses apps from CSV using csvutil
func (am *AppModel) ParseApps() ([]App, error) {
	var apps []App
	records, err := utils.ReadCSV(am.config.CSVFilePath)
	if err != nil {
		return nil, err
	}

	// Convert records to CSV format
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err = writer.WriteAll(records)
	if err != nil {
		return nil, err
	}
	writer.Flush()

	// Unmarshal CSV into struct
	if err := csvutil.Unmarshal(buf.Bytes(), &apps); err != nil {
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

// GetAllApps: Returns apps with pagination and filters
func (am *AppModel) GetAllApps(limit int, page int, priceFilter string) ([]string, error) {
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
				return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid price value")
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

// func parseDate(dateStr string) string {
// 	layout := "January 2, 2006"
// 	parsedTime, _ := time.Parse(layout, dateStr)
// 	return parsedTime.Format("2006-01-02")
// }
